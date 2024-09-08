package sourcer

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"
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
	Source     Sourcer
	shutdownCh chan<- struct{}
}

// ReadFn reads the data from the source.
func (fs *Service) ReadFn(stream sourcepb.Source_ReadFnServer) error {
	ctx := stream.Context()
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// handle panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
				fs.shutdownCh <- struct{}{}
				errCh <- fmt.Errorf("panic: %v", r)
			}
		}()

		for {
			// Receive read requests from the stream
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("error receiving from stream: %v", err)
				errCh <- err
				return
			}

			messageCh := make(chan Message)
			go func() {
				defer close(messageCh)
				request := readRequest{
					count:   req.Request.GetNumRecords(),
					timeout: time.Duration(req.Request.GetTimeoutInMs()) * time.Millisecond,
				}
				fs.Source.Read(ctx, &request, messageCh)
			}()

			// Read messages from the channel and send them to the stream.
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
					Status: &sourcepb.ReadResponse_Status{
						Eot:  false,
						Code: 0,
					},
				}
				// The error here is returned by the stream, which is already a gRPC error
				if err := stream.Send(element); err != nil {
					errCh <- err
					return
				}
			}
			err = stream.Send(&sourcepb.ReadResponse{
				Status: &sourcepb.ReadResponse_Status{
					Eot:  true,
					Code: 0,
				},
			})
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Wait for the goroutine to finish
	wg.Wait()
	// Check if there was any error in the goroutine
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// ackRequest implements the AckRequest interface and is used in the ack handler.
type ackRequest struct {
	offset Offset
}

// Offset returns the offset of the record to ack.
func (a *ackRequest) Offset() Offset {
	return a.offset
}

// AckFn acknowledges the data from the source.
func (fs *Service) AckFn(stream sourcepb.Source_AckFnServer) error {
	ctx := stream.Context()

	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

	for {
		// Receive ack requests from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&sourcepb.AckResponse{
				Result: &sourcepb.AckResponse_Result{},
			})
		}
		if err != nil {
			log.Printf("error receiving from stream: %v", err)
			return err
		}

		request := ackRequest{
			offset: NewOffset(req.Request.Offset.GetOffset(), req.Request.Offset.GetPartitionId()),
		}
		fs.Source.Ack(ctx, &request)
	}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sourcepb.ReadyResponse, error) {
	return &sourcepb.ReadyResponse{Ready: true}, nil
}

// PendingFn returns the number of pending messages.
func (fs *Service) PendingFn(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PendingResponse, error) {
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside sourcer handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

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

func (fs *Service) PartitionsFn(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PartitionsResponse, error) {
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

	partitions := fs.Source.Partitions(ctx)
	return &sourcepb.PartitionsResponse{
		Result: &sourcepb.PartitionsResponse_Result{
			Partitions: partitions,
		},
	}, nil
}
