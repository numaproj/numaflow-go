package mapstreamer

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	mapstreampb "github.com/numaproj/numaflow-go/pkg/apis/proto/mapstream/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/mapstream.sock"
	serverInfoFilePath    = "/var/run/numaflow/mapstreamer-server-info"
)

// Service implements the proto gen server interface and contains the map
// streaming function.
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
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders())
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
				Result: &mapstreampb.MapStreamResponse_Result{
					Keys:  message.Keys(),
					Value: message.Value(),
					Tags:  message.Tags(),
				},
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
