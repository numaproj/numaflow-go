package sourcer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"time"

	"golang.org/x/sync/errgroup"
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
			log.Printf("error processing requests: %v", err)
			fs.shutdownCh <- struct{}{}
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

// recvWithContext wraps stream.Recv() to respect context cancellation for ReadFn.
func recvWithContextRead(ctx context.Context, stream sourcepb.Source_ReadFnServer) (*sourcepb.ReadRequest, error) {
	type recvResult struct {
		req *sourcepb.ReadRequest
		err error
	}

	resultCh := make(chan recvResult, 1)
	go func() {
		req, err := stream.Recv()
		resultCh <- recvResult{req: req, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return result.req, result.err
	}
}

// receiveReadRequests receives read requests from the client and invokes the source Read method.
// writes the read data to the message channel.
func (fs *Service) receiveReadRequests(ctx context.Context, stream sourcepb.Source_ReadFnServer) error {
	messageCh := make(chan Message)
	eg, groupCtx := errgroup.WithContext(ctx)
	req, err := recvWithContextRead(groupCtx, stream)
	if err == io.EOF {
		log.Printf("end of read stream")
		return err
	}
	if err != nil {
		log.Printf("error receiving from read stream: %v", err)
		return err
	}

	eg.Go(func() (err error) {
		// handle panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
				err = fmt.Errorf("panic inside source handler: %v", r)
				return
			}
			close(messageCh)
		}()
		request := readRequest{
			count:   req.Request.GetNumRecords(),
			timeout: time.Duration(req.Request.GetTimeoutInMs()) * time.Millisecond,
		}
		fs.Source.Read(ctx, &request, messageCh)
		return nil
	})

	// invoke the processReadData method to send the read data to the client.
	eg.Go(func() error {
		return fs.processReadData(groupCtx, stream, messageCh)
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// processReadData processes the read data and sends it to the client.
func (fs *Service) processReadData(ctx context.Context, stream sourcepb.Source_ReadFnServer, messageCh <-chan Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg, ok := <-messageCh:
		if !ok {
			break
		}
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
	offsets []Offset
}

// Offsets returns the offsets to be acknowledged.
func (a *ackRequest) Offsets() []Offset {
	return a.offsets
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

// recvWithContext wraps stream.Recv() to respect context cancellation for AckFn.
func recvWithContextAck(ctx context.Context, stream sourcepb.Source_AckFnServer) (*sourcepb.AckRequest, error) {
	type recvResult struct {
		req *sourcepb.AckRequest
		err error
	}

	resultCh := make(chan recvResult, 1)
	go func() {
		req, err := stream.Recv()
		resultCh <- recvResult{req: req, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return result.req, result.err
	}
}

// receiveAckRequests receives ack requests from the client and invokes the source Ack method.
func (fs *Service) receiveAckRequests(ctx context.Context, stream sourcepb.Source_AckFnServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside source handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
			err = fmt.Errorf("panic inside source handler: %v", r)
		}
	}()

	req, err := recvWithContextAck(ctx, stream)
	if err == io.EOF {
		log.Printf("end of ack stream")
		return err
	}
	if err != nil {
		log.Printf("error receiving from ack stream: %v", err)
		return err
	}

	offsets := make([]Offset, len(req.Request.GetOffsets()))
	for i, offset := range req.Request.GetOffsets() {
		offsets[i] = NewOffset(offset.GetOffset(), offset.GetPartitionId())
	}

	request := ackRequest{
		offsets: offsets,
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
