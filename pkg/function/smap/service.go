package smap

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/smapfn"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service implements the proto gen server interface and contains the map
// streaming handler
type Service struct {
	smapfn.UnimplementedMapStreamServer

	MapperStream function.MapStreamHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*smapfn.ReadyResponse, error) {
	return &smapfn.ReadyResponse{Ready: true}, nil
}

// MapStreamFn applies a function to each datum element and returns a stream.
func (fs *Service) MapStreamFn(d *smapfn.MapStreamRequest, stream smapfn.MapStream_MapStreamFnServer) error {
	var hd = function.NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime(), function.NewHandlerDatumMetadata("id", 1))
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
			element := &smapfn.MapStreamResponse{
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
