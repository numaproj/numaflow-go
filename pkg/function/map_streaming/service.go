package map_streaming

import (
	"context"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/map_streaming/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service implements the proto gen server interface and contains the map operation
// handler and the reduce operation handler.
type Service struct {
	v1.UnimplementedMapStreamServer

	MapperStream function.MapStreamHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

// MapStreamFn applies a function to each datum element and returns a stream.
func (fs *Service) MapStreamFn(d *v1.MapStreamRequest, stream v1.MapStream_MapStreamFnServer) error {
	var hd = function.NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime(), nil)
	ctx := stream.Context()
	messageCh := make(chan function.Message)

	done := make(chan bool)
	go func() {
		fs.MapperStream.HandleDo(ctx, d.GetKeys(), hd, messageCh)
		done <- true
	}()
	finished := false
	for {
		select {
		case <-done:
			finished = true
		case message, ok := <-messageCh:
			if !ok {
				// Channel already closed, not closing again.
				return nil
			}
			element := &v1.MapStreamResponse{
				Keys:  message.Keys(),
				Value: message.Value(),
				Tags:  message.Tags(),
			}
			err := stream.Send(element)
			// the error here is returned by stream.Send() which is already a gRPC error
			if err != nil {
				// Channel may or may not be closed, as we are not sure leave it to GC.
				return err
			}
		default:
			if finished {
				close(messageCh)
				return nil
			}
		}
	}
}
