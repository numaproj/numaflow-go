package mapper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"golang.org/x/sync/errgroup"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/map.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
	EnvUDContainerType    = "NUMAFLOW_UD_CONTAINER_TYPE"
)

var errMapHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", os.Getenv(EnvUDContainerType))

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	mappb.UnimplementedMapServer
	Mapper     Mapper
	shutdownCh chan<- struct{}
	once       sync.Once
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// recvWithContext wraps stream.Recv() to respect context cancellation.
func recvWithContext(ctx context.Context, stream mappb.Map_MapFnServer) (*mappb.MapRequest, error) {
	type recvResult struct {
		req *mappb.MapRequest
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

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(stream mappb.Map_MapFnServer) error {
	// perform handshake with client before processing requests
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Use error group to manage goroutines, the groupCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the wait returns.
	g, groupCtx := errgroup.WithContext(ctx)

	// Channel to collect responses
	responseCh := make(chan *mappb.MapResponse, 500) // FIXME: identify the right buffer size
	defer close(responseCh)

	// Dedicated goroutine to send responses to the stream
	g.Go(func() error {
		for {
			select {
			case resp := <-responseCh:
				if err := stream.Send(resp); err != nil {
					log.Printf("Failed to send response: %v", err)
					return err
				}
			case <-groupCtx.Done():
				return groupCtx.Err()
			}
		}
	})

	var readErr error
	// Read requests from the stream and process them
outer:
	for {
		req, err := recvWithContext(groupCtx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the MapFn")
			break outer
		}
		if errors.Is(err, io.EOF) {
			log.Printf("EOF received, stopping the MapFn")
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
		g.Go(func() error {
			return fs.handleRequest(groupCtx, req, responseCh)
		})
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the MapFn with err, %s", err)
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
func (fs *Service) performHandshake(stream mappb.Map_MapFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive handshake: %v", err)
	}
	if req.GetHandshake() == nil || !req.GetHandshake().GetSot() {
		return status.Errorf(codes.InvalidArgument, "invalid handshake")
	}
	handshakeResponse := &mappb.MapResponse{
		Handshake: &mappb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return fmt.Errorf("sending handshake response to client over gRPC stream: %w", err)
	}
	return nil
}

// handleRequest processes each request and sends the response to the response channel.
func (fs *Service) handleRequest(ctx context.Context, req *mappb.MapRequest, responseCh chan<- *mappb.MapResponse) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside map handler: %v %v", r, string(debug.Stack()))
			st, _ := status.Newf(codes.Internal, "%s: %v", errMapHandlerPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
		}
	}()

	request := req.GetRequest()
	hd := NewHandlerDatum(request.GetValue(), request.GetEventTime().AsTime(), request.GetWatermark().AsTime(), request.GetHeaders())
	messages := fs.Mapper.Map(ctx, request.GetKeys(), hd)
	var elements []*mappb.MapResponse_Result
	for _, m := range messages.Items() {
		elements = append(elements, &mappb.MapResponse_Result{
			Keys:  m.Keys(),
			Value: m.Value(),
			Tags:  m.Tags(),
		})
	}
	resp := &mappb.MapResponse{
		Results: elements,
		Id:      req.GetId(),
	}
	select {
	case responseCh <- resp:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
