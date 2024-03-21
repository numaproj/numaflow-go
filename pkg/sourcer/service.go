package sourcer

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64 // 64MB
	address               = "/var/run/numaflow/source.sock"
	serverInfoFilePath    = "/var/run/numaflow/sourcer-server-info"
)

// Service implements the proto gen server interface
type Service struct {
	sourcepb.UnimplementedSourceServer
	Source Sourcer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sourcepb.ReadyResponse, error) {
	return &sourcepb.ReadyResponse{Ready: true}, nil
}

// PendingFn returns the number of pending messages.
func (fs *Service) PendingFn(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PendingResponse, error) {
	return &sourcepb.PendingResponse{Result: &sourcepb.PendingResponse_Result{
		Count: fs.Source.Pending(ctx),
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

// ReadFn reads the data from the source.
func (fs *Service) ReadFn(d *sourcepb.ReadRequest, stream sourcepb.Source_ReadFnServer) error {
	request := readRequest{
		count:   d.Request.GetNumRecords(),
		timeout: time.Duration(d.Request.GetTimeoutInMs()) * time.Millisecond,
	}
	ctx := stream.Context()
	messageCh := make(chan Message)

	// Start the read in a goroutine
	go func() {
		defer close(messageCh)
		fs.Source.Read(ctx, &request, messageCh)
	}()

	// Read messages from the channel and send them to the stream, until the channel is closed
	for msg := range messageCh {
		offset := &sourcepb.Offset{
			Offset:      msg.Offset().Value(),
			PartitionId: msg.Offset().PartitionId(),
		}
		element := &sourcepb.ReadResponse{
			Result: &sourcepb.ReadResponse_Result{
				Payload:   msg.Value(),
				Offset:    offset,
				EventTime: timestamppb.New(msg.EventTime()),
				Keys:      msg.Keys(),
				Headers:   msg.Headers(),
			},
		}
		// The error here is returned by the stream, which is already a gRPC error
		if err := stream.Send(element); err != nil {
			// The channel may or may not be closed, as we are not sure, we leave it to GC.
			return err
		}
	}
	return nil
}

// ackRequest implements the AckRequest interface and is used in the ack handler.
type ackRequest struct {
	offsets []Offset
}

// Offsets returns the offsets of the records to ack.
func (a *ackRequest) Offsets() []Offset {
	return a.offsets
}

// AckFn applies a function to each datum element.
func (fs *Service) AckFn(ctx context.Context, d *sourcepb.AckRequest) (*sourcepb.AckResponse, error) {
	offsets := make([]Offset, len(d.Request.GetOffsets()))
	for i, offset := range d.Request.GetOffsets() {
		offsets[i] = NewOffset(offset.GetOffset(), offset.GetPartitionId())
	}

	request := ackRequest{
		offsets: offsets,
	}
	fs.Source.Ack(ctx, &request)
	return &sourcepb.AckResponse{
		Result: &sourcepb.AckResponse_Result{},
	}, nil
}

func (fs *Service) PartitionsFn(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PartitionsResponse, error) {
	partitions := fs.Source.Partitions(ctx)
	return &sourcepb.PartitionsResponse{
		Result: &sourcepb.PartitionsResponse_Result{
			Partitions: partitions,
		},
	}, nil
}
