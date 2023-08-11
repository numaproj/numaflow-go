package source

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/source/model"
)

// readRequest implements the ReadRequest interface and is used in the read handler.
type readRequest struct {
	count uint64
}

func (r *readRequest) Count() uint64 {
	return r.count
}

// Service implements the proto gen server interface
type Service struct {
	sourcepb.UnimplementedSourceServer
	PendingHandler PendingHandler
	ReadHandler    ReadHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sourcepb.ReadyResponse, error) {
	return &sourcepb.ReadyResponse{Ready: true}, nil
}

// Pending returns the number of pending messages.
func (fs *Service) Pending(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PendingResponse, error) {
	return &sourcepb.PendingResponse{Result: &sourcepb.PendingResponse_Result{
		Count: fs.PendingHandler.HandleDo(ctx),
	}}, nil
}

// Read reads the data from the source.
func (fs *Service) Read(d *sourcepb.ReadRequest, stream sourcepb.Source_ReadFnServer) error {
	request := readRequest{
		count: d.Request.GetNumRecords(),
	}
	ctx := stream.Context()
	messageCh := make(chan model.Message)

	done := make(chan bool)
	go func() {
		fs.ReadHandler.HandleDo(ctx, &request, messageCh)
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
			offset := &sourcepb.Offset{
				Offset:      message.Offset().Value(),
				PartitionId: message.Offset().PartitionId(),
			}
			element := &sourcepb.DatumResponse{
				Result: &sourcepb.DatumResponse_Result{
					Payload:   message.Value(),
					Offset:    offset,
					EventTime: timestamppb.New(message.EventTime()),
					Keys:      message.Keys(),
				},
			}
			err := stream.Send(element)
			// stream returns the error here.Send() which is already a gRPC error
			if err != nil {
				// The channel may or may not be closed, as we are not sure, leave it to GC.
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
