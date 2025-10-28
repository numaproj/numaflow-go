package mapstreamer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"

	"golang.org/x/sync/errgroup"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/mapstream.sock"
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

var errMapStreamHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// Service implements the proto gen server interface and contains the map
// streaming function.
type Service struct {
	mappb.UnimplementedMapServer
	shutdownCh   chan<- struct{}
	MapperStream MapStreamer
	once         sync.Once
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

// MapFn applies a function to each request element and streams the results back.
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
			messageCh := make(chan Message)
			workerGroup, innerCtx := errgroup.WithContext(groupCtx)
			workerGroup.Go(func() error {
				return fs.invokeHandler(innerCtx, req, messageCh)
			})
			workerGroup.Go(func() error {
				return fs.writeResponseToClient(innerCtx, stream, req.GetId(), messageCh)
			})
			return workerGroup.Wait()
		})
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		fs.once.Do(func() {
			log.Printf("Stopping the MapFn with err, %s", err)
			select {
			case fs.shutdownCh <- struct{}{}:
				// signal enqueued
			default:
				log.Println("Shutdown signal already enqueued or watcher exited; skipping shutdown send")
			}
		})
		return err
	}

	// check if there was an error while reading from the stream
	if readErr != nil {
		return status.Errorf(codes.Internal, "%s", readErr.Error())
	}

	return nil
}

// invokeHandler handles the map stream invocation.
func (fs *Service) invokeHandler(ctx context.Context, req *mappb.MapRequest, messageCh chan<- Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside mapStream handler: %v %v", r, string(debug.Stack()))
			st, _ := status.Newf(codes.Internal, "%s: %v", errMapStreamHandlerPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
			return
		}
	}()
	streamReq := req.GetRequest()
	hd := NewHandlerDatum(streamReq.GetValue(), streamReq.GetEventTime().AsTime(), streamReq.GetWatermark().AsTime(), streamReq.GetHeaders())
	fs.MapperStream.MapStream(ctx, req.GetRequest().GetKeys(), hd, messageCh)
	return nil
}

// writeResponseToClient writes the response back to the client.
func (fs *Service) writeResponseToClient(ctx context.Context, stream mappb.Map_MapFnServer, reqID string, messageCh <-chan Message) error {
	for {
		select {
		case message, ok := <-messageCh:
			if !ok {
				// Send EOT message since we are done processing the request.
				eotMessage := &mappb.MapResponse{
					Status: &mappb.TransmissionStatus{Eot: true},
					Id:     reqID,
				}
				if err := stream.Send(eotMessage); err != nil {
					return err
				}
				return nil
			}
			element := &mappb.MapResponse{
				Results: []*mappb.MapResponse_Result{
					{
						Keys:  message.Keys(),
						Value: message.Value(),
						Tags:  message.Tags(),
					},
				},
				Id: reqID,
			}
			if err := stream.Send(element); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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
