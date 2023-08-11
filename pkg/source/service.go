package source

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/source/model"
)

// Service implements the proto gen server interface
type Service struct {
	sourcepb.UnimplementedSourceServer
	PendingHandler PendingHandler
	ReadHandler    ReadHandler
	AckHandler     AckHandler
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

// readRequest implements the ReadRequest interface and is used in the read handler.
type readRequest struct {
	count   uint64
	timeout time.Duration
}

func (r *readRequest) TimeOut() time.Duration {
	return r.timeout
}

func (r *readRequest) Count() uint64 {
	return r.count
}

// Read reads the data from the source.
func (fs *Service) Read(d *sourcepb.ReadRequest, stream sourcepb.Source_ReadFnServer) error {
	request := readRequest{
		count:   d.Request.GetNumRecords(),
		timeout: time.Duration(d.Request.GetTimeoutInMs()) * time.Millisecond,
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

// ackRequest implements the AckRequest interface and is used in the ack handler.
type ackRequest struct {
	offsets []model.Offset
}

// Offsets returns the offsets of the records to ack.
func (a *ackRequest) Offsets() []model.Offset {
	return a.offsets
}

// AckFn applies a function to each datum element.
func (fs *Service) AckFn(ctx context.Context, d *sourcepb.AckRequest) (*sourcepb.AckResponse, error) {
	offsets := make([]model.Offset, len(d.Request.GetOffsets()))
	for i, offset := range d.Request.GetOffsets() {
		offsets[i] = model.NewOffset(offset.GetOffset(), offset.GetPartitionId())
	}

	request := ackRequest{
		offsets: offsets,
	}
	fs.AckHandler.HandleDo(ctx, &request)
	return &sourcepb.AckResponse{
		Result: &sourcepb.AckResponse_Result{},
	}, nil
}
