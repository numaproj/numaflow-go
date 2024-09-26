package sourcetransformer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1"
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
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

// SourceTransformFn applies a function to each request element.
// In addition to map function, SourceTransformFn also supports assigning a new event time to response.
// SourceTransformFn can be used only at source vertex by source data transformer.
func (fs *Service) SourceTransformFn(stream v1.SourceTransform_SourceTransformFnServer) error {
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside sourcetransform handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("reading handshake message from stream: %w", err)
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("invalid handshake message: %+v", req)
	}

	handshakeResponse := &v1.SourceTransformResponse{
		Handshake: &v1.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return fmt.Errorf("sending handshake response to client over gRPC stream: %w", err)
	}

	ctx := stream.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	grp, grpCtx := errgroup.WithContext(ctx)

	senderCh := make(chan *v1.SourceTransformResponse, 500) // TODO: identify the right buffer size
	// goroutine to send the response to the stream
	grp.Go(func() error {
		for {
			select {
			case <-grpCtx.Done():
				return grpCtx.Err()
			default:
			}
			if err := stream.Send(<-senderCh); err != nil {
				cancel()
				return err
			}
		}
	})

	for {
		d, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if d.Handshake != nil {
			return fmt.Errorf("expected source transform messages, received handshake message")
		}

		req := d.Request
		grp.Go(func() error {
			var hd = NewHandlerDatum(req.GetValue(), req.EventTime.AsTime(), req.Watermark.AsTime(), req.Headers)
			messageTs := fs.Transformer.Transform(grpCtx, req.GetKeys(), hd)
			var results []*v1.SourceTransformResponse_Result
			for _, m := range messageTs.Items() {
				results = append(results, &v1.SourceTransformResponse_Result{
					EventTime: timestamppb.New(m.EventTime()),
					Keys:      m.Keys(),
					Value:     m.Value(),
					Tags:      m.Tags(),
				})
			}
			resp := &v1.SourceTransformResponse{
				Results: results,
				Id:      req.GetId(),
			}
			select {
			case senderCh <- resp:
			case <-grpCtx.Done():
				return grpCtx.Err()
			}
			return nil
		})
	}

	if err := grp.Wait(); err != nil {
		statusErr := status.Errorf(codes.Internal, err.Error())
		return statusErr
	}
	return nil
}
