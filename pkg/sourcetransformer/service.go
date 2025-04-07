package sourcetransformer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/sourcetransform.sock"
	serverInfoFilePath    = "/var/run/numaflow/sourcetransformer-server-info"
)

// Service implements the proto gen server interface and contains the transformer operation
// handler.
type Service struct {
	v1.UnimplementedSourceTransformServer
	Transformer SourceTransformer
	shutdownCh  chan<- struct{}
	once        sync.Once
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

var errTransformerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// recvWithContext wraps stream.Recv() to respect context cancellation.
func recvWithContext(ctx context.Context, stream v1.SourceTransform_SourceTransformFnServer) (*v1.SourceTransformRequest, error) {
	type recvResult struct {
		req *v1.SourceTransformRequest
		err error
	}

	resultCh := make(chan recvResult, 1)
	go func() {
		req, err := stream.Recv()
		resultCh <- recvResult{req: req, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return result.req, result.err
	}
}

// SourceTransformFn applies a function to each request element.
// In addition to map function, SourceTransformFn also supports assigning a new event time to response.
// SourceTransformFn can be used only at source vertex by source data transformer.
func (fs *Service) SourceTransformFn(stream v1.SourceTransform_SourceTransformFnServer) error {

	// perform handshake with client before processing requests
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Use error group to manage goroutines, the groupCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the wait returns.
	grp, groupCtx := errgroup.WithContext(ctx)

	senderCh := make(chan *v1.SourceTransformResponse, 500) // FIXME: identify the right buffer size
	// goroutine to send the responses back to the client
	grp.Go(func() error {
		for {
			select {
			case <-groupCtx.Done():
				return groupCtx.Err()
			case resp := <-senderCh:
				if err := stream.Send(resp); err != nil {
					log.Printf("Failed to send response: %v", err)
					return fmt.Errorf("failed to send response to client: %w", err)
				}
			}
		}
	})

	var readErr error
outer:
	for {
		d, err := recvWithContext(groupCtx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the SourceTransformFn")
			break outer
		}
		if errors.Is(err, io.EOF) {
			log.Printf("EOF received, stopping the SourceTransformFn")
			break outer
		}
		if err != nil {
			log.Printf("Failed to receive request: %v", err)
			readErr = err
			// read loop is not part of the error group, so we need to cancel the context
			// to signal the other goroutines to stop processing.
			cancel()
			break outer
		}
		grp.Go(func() (err error) {
			return fs.handleRequest(groupCtx, d, senderCh)
		})
	}

	// wait for all the goroutines to finish, if any of the goroutines return an error, wait will return that error immediately.
	if err := grp.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the SourceTransformFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
		})
		return err
	}

	// check if there was an error while reading from the stream
	if readErr != nil {
		return status.Errorf(codes.Internal, "%s", readErr.Error())
	}

	return nil
}

// performHandshake handles the handshake logic at the start of the stream.
func (fs *Service) performHandshake(stream v1.SourceTransform_SourceTransformFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive handshake: %v", err)
	}
	if req.GetHandshake() == nil || !req.GetHandshake().GetSot() {
		return status.Errorf(codes.InvalidArgument, "invalid handshake")
	}
	handshakeResponse := &v1.SourceTransformResponse{
		Handshake: &v1.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return fmt.Errorf("sending handshake response to client over gRPC stream: %w", err)
	}
	return nil
}

// handleRequest processes each request and sends the response to the response channel.
func (fs *Service) handleRequest(ctx context.Context, req *v1.SourceTransformRequest, responseCh chan<- *v1.SourceTransformResponse) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside handler: %v %v", r, string(debug.Stack()))
			st, _ := status.Newf(codes.Internal, "%s: %v", errTransformerPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
		}
	}()

	request := req.GetRequest()
	hd := NewHandlerDatum(request.GetValue(), request.GetEventTime().AsTime(), request.GetWatermark().AsTime(), request.GetHeaders())
	messages := fs.Transformer.Transform(ctx, request.GetKeys(), hd)
	var elements []*v1.SourceTransformResponse_Result
	for _, m := range messages.Items() {
		elements = append(elements, &v1.SourceTransformResponse_Result{
			Keys:      m.Keys(),
			Value:     m.Value(),
			Tags:      m.Tags(),
			EventTime: timestamppb.New(m.EventTime()),
		})
	}
	resp := &v1.SourceTransformResponse{
		Results: elements,
		Id:      req.GetRequest().GetId(),
	}
	select {
	case responseCh <- resp:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
