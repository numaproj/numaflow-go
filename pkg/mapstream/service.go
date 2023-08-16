package mapstream

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	mapstreampb "github.com/numaproj/numaflow-go/pkg/apis/proto/mapstream/v1"
)

// Service implements the proto gen server interface and contains the map
// streaming handler
type Service struct {
	mapstreampb.UnimplementedMapStreamServer

	MapperStream MapStreamer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mapstreampb.ReadyResponse, error) {
	return &mapstreampb.ReadyResponse{Ready: true}, nil
}

// MapStreamFn applies a function to each request element and streams the results back.
func (fs *Service) MapStreamFn(d *mapstreampb.MapStreamRequest, stream mapstreampb.MapStream_MapStreamFnServer) error {
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime())
	ctx := stream.Context()
	messageCh := make(chan Message)

	done := make(chan bool)
	go func() {
		fs.MapperStream.MapStream(ctx, d.GetKeys(), hd, messageCh)
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
			element := &mapstreampb.MapStreamResponse{
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
