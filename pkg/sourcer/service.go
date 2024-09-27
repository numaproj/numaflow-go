package sourcer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
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

	if err := fs.performReadHandshake(stream); err != nil {
		return err
	}

	for {
		if err := fs.receiveReadRequests(ctx, stream); err != nil {
			// If the error is EOF, it means the stream has been closed.
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

// performReadHandshake performs the handshake with the client before starting the read process.
func (fs *Service) performReadHandshake(stream sourcepb.Source_ReadFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("error receiving handshake from stream: %v", err)
		return err
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("expected handshake message")
	}

	handshakeResponse := &sourcepb.ReadResponse{
		Status: &sourcepb.ReadResponse_Status{
			Eot:  false,
			Code: sourcepb.ReadResponse_Status_SUCCESS,
		},
		Handshake: &sourcepb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return err
	}

	return nil
}

// receiveReadRequests receives read requests from the client and invokes the source Read method.
// writes the read data to the message channel.
func (fs *Service) receiveReadRequests(ctx context.Context, stream sourcepb.Source_ReadFnServer) error {
	messageCh := make(chan Message)
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

	req, err := stream.Recv()
	if err == io.EOF {
		log.Printf("end of read stream")
		return err
	}
	if err != nil {
		log.Printf("error receiving from read stream: %v", err)
		return err
	}

	go func() {
		defer close(messageCh)
		request := readRequest{
			count:   req.Request.GetNumRecords(),
			timeout: time.Duration(req.Request.GetTimeoutInMs()) * time.Millisecond,
		}
		fs.Source.Read(ctx, &request, messageCh)
	}()

	// invoke the processReadData method to send the read data to the client.
	return fs.processReadData(stream, messageCh)
}

// processReadData processes the read data and sends it to the client.
func (fs *Service) processReadData(stream sourcepb.Source_ReadFnServer, messageCh <-chan Message) error {
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
		if err := stream.Send(element); err != nil {
			return err
		}
	}
	err := stream.Send(&sourcepb.ReadResponse{
		Status: &sourcepb.ReadResponse_Status{
			Eot:  true,
			Code: 0,
		},
	})
	if err != nil {
		return err
	}
	return nil
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

	if err := fs.performAckHandshake(stream); err != nil {
		return err
	}

	for {
		err := fs.receiveAckRequests(ctx, stream)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

// performAckHandshake performs the handshake with the client before starting the ack process.
func (fs *Service) performAckHandshake(stream sourcepb.Source_AckFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("error receiving handshake from stream: %v", err)
		return err
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("expected handshake message")
	}

	handshakeResponse := &sourcepb.AckResponse{
		Result: &sourcepb.AckResponse_Result{
			Success: &emptypb.Empty{},
		},
		Handshake: &sourcepb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return err
	}

	return nil
}

// receiveAckRequests receives ack requests from the client and invokes the source Ack method.
func (fs *Service) receiveAckRequests(ctx context.Context, stream sourcepb.Source_AckFnServer) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
		}
	}()

	req, err := stream.Recv()
	if err == io.EOF {
		log.Printf("end of ack stream")
		return err
	}
	if err != nil {
		log.Printf("error receiving from ack stream: %v", err)
		return err
	}

	request := ackRequest{
		offset: NewOffset(req.Request.Offset.GetOffset(), req.Request.Offset.GetPartitionId()),
	}
	fs.Source.Ack(ctx, &request)

	ackResponse := &sourcepb.AckResponse{
		Result: &sourcepb.AckResponse_Result{
			Success: &emptypb.Empty{},
		},
	}
	if err := stream.Send(ackResponse); err != nil {
		log.Printf("error sending ack response: %v", err)
		return err
	}
	return nil
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
