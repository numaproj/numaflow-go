package reducestreamer

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	reducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/reducestream.sock"
	winStartTime          = "x-numaflow-win-start-time"
	winEndTime            = "x-numaflow-win-end-time"
	delimiter             = ":"
	serverInfoFilePath    = "/var/run/numaflow/reducestreamer-server-info"
)

// Service implements the proto gen server interface and contains the reduceStream operation handler.
type Service struct {
	reducepb.UnimplementedReduceServer
	creatorHandle ReduceStreamerCreator
	shutdownCh    chan<- struct{}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*reducepb.ReadyResponse, error) {
	return &reducepb.ReadyResponse{Ready: true}, nil
}

// ReduceFn applies a reduce function to a request stream and streams the results.
func (fs *Service) ReduceFn(stream reducepb.Reduce_ReduceFnServer) error {
	var (
		err error
		ctx = stream.Context()
		g   errgroup.Group
	)

	taskManager := newReduceTaskManager(fs.creatorHandle, fs.shutdownCh)

	// err group for the go routine which reads from the output channel and sends to the stream
	g.Go(func() error {
		for output := range taskManager.OutputChannel() {
			sendErr := stream.Send(output)
			if sendErr != nil {
				return sendErr
			}
		}
		return nil
	})

	// Start a goroutine to receive messages from the gRPC stream and forward them to recvCh.
	// Any error encountered (including io.EOF) is sent to recvErrCh.
	recvCh := make(chan *reducepb.ReduceRequest, 1)
	recvErrCh := make(chan error, 1)
	go func() {
		for {
			d, err := stream.Recv()
			if err != nil {
				recvErrCh <- err
				return
			}
			recvCh <- d
		}
	}()

	// Main loop to process incoming requests and handle errors.
	// Uses select to:
	// - Process messages from the stream (via recvCh)
	// - Handle errors from task goroutines (via taskManager.ErrorChannel())
	// - Handle stream errors and EOF (via recvErrCh)
	// This ensures the main loop remains responsive to both stream and task errors at all times.

	// read messages from the stream and write the messages to corresponding channels
	// if the channel is not created, create the channel and invoke the reduceFn
loop:
	for {
		select {
		case errFromTask := <-taskManager.ErrorChannel():
			fs.shutdownCh <- struct{}{}
			return errFromTask
		case d := <-recvCh:
			switch d.Operation.Event {
			case reducepb.ReduceRequest_WindowOperation_OPEN:
				// create a new reduce task and start the reduce operation
				err = taskManager.CreateTask(ctx, d)
				if err != nil {
					statusErr := status.Errorf(codes.Internal, "%s", err.Error())
					return statusErr
				}
			case reducepb.ReduceRequest_WindowOperation_APPEND:
				// append the datum to the reduce task
				err = taskManager.AppendToTask(ctx, d)
				if err != nil {
					statusErr := status.Errorf(codes.Internal, "%s", err.Error())
					return statusErr
				}
			}
		case recvErr := <-recvErrCh:
			if recvErr == io.EOF {
				taskManager.CloseAll()
				break loop
			}
			return recvErr
		}
	}

	taskManager.WaitAll()
	// wait for the go routine which reads from the output channel and sends to the stream to return
	err = g.Wait()
	if err != nil {
		statusErr := status.Errorf(codes.Internal, "%s", err.Error())
		return statusErr
	}

	return nil
}
