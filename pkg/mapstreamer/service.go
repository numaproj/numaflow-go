package mapstreamer

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
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/mapstream.sock"
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

// Service implements the proto gen server interface and contains the map
// streaming function.
type Service struct {
	mappb.UnimplementedMapServer
	shutdownCh   chan<- struct{}
	MapperStream MapStreamer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a function to each request element and streams the results back.
func (fs *Service) MapFn(stream mappb.Map_MapFnServer) error {
	// perform handshake with client before processing requests
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to receive request: %v", err)
			return err
		}

		messageCh := make(chan Message)
		g, groupCtx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return fs.invokeHandler(groupCtx, req, messageCh)
		})

		g.Go(func() error {
			return fs.writeResponseToClient(groupCtx, stream, req.GetId(), messageCh)
		})

		// Wait for the error group to finish
		if err := g.Wait(); err != nil {
			log.Printf("error processing requests: %v", err)
			if err == io.EOF {
				return nil
			}
			fs.shutdownCh <- struct{}{}
			return status.Errorf(codes.Internal, "error processing requests: %v", err)
		}
	}

	return nil
}

// invokeHandler handles the map stream invocation.
func (fs *Service) invokeHandler(ctx context.Context, req *mappb.MapRequest, messageCh chan<- Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside mapStream handler: %v %v", r, string(debug.Stack()))
			err = fmt.Errorf("panic inside mapStream handler: %v", r)
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
