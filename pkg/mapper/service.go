package mapper

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
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
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	mappb.UnimplementedMapServer
	Mapper     Mapper
	shutdownCh chan<- struct{}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(stream mappb.Map_MapFnServer) error {
	// perform handshake with client before processing requests
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Use error group to manage goroutines, the grpCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the wait returns.
	g, grpCtx := errgroup.WithContext(ctx)

	// Channel to collect responses
	responseCh := make(chan *mappb.MapResponse, 500) // FIXME: identify the right buffer size
	defer close(responseCh)

	// Dedicated goroutine to send responses to the stream
	g.Go(func() error {
		for {
			select {
			case resp := <-responseCh:
				if err := stream.Send(resp); err != nil {
					log.Printf("failed to send response: %v", err)
					return err
				}
			case <-grpCtx.Done():
				return nil
			}
		}
	})

	var readErr error
	// Read requests from the stream and process them
outer:
	for {
		select {
		case <-grpCtx.Done():
			break outer
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				break outer
			}
			if err != nil {
				log.Printf("failed to receive request: %v", err)
				readErr = err
				// read loop is not part of the error group, so we need to cancel the context
				// to signal the other goroutines to stop processing.
				cancel()
				break outer
			}
			g.Go(func() error {
				return fs.handleRequest(grpCtx, req, responseCh)
			})
		}
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fs.shutdownCh <- struct{}{}
		return status.Errorf(codes.Internal, "error processing requests: %v", err)
	}

	// check if there was an error while reading from the stream
	if readErr != nil {
		return status.Errorf(codes.Internal, readErr.Error())
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
func (fs *Service) handleRequest(ctx context.Context, req *mappb.MapRequest, responseCh chan<- *mappb.MapResponse) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside map handler: %v %v", r, string(debug.Stack()))
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
		return nil
	}
	return nil
}
