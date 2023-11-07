package reducer

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
	address               = "/var/run/numaflow/reduce.sock"
	winStartTime          = "x-numaflow-win-start-time"
	winEndTime            = "x-numaflow-win-end-time"
	delimiter             = ":"
)

// Service implements the proto gen server interface and contains the reduce operation handler.
type Service struct {
	reducepb.UnimplementedReduceServer
	Reducer Reducer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*reducepb.ReadyResponse, error) {
	return &reducepb.ReadyResponse{Ready: true}, nil
}

// ReduceFn applies a reduce function to a request stream and returns a list of results.
func (fs *Service) ReduceFn(stream reducepb.Reduce_ReduceFnServer) error {
	var (
		err error
		ctx = stream.Context()
		g   errgroup.Group
	)

	taskManager := newReduceTaskManager()

	// err group for the go routine which reads from the output channel and sends to the stream
	g.Go(func() error {
		for output := range taskManager.OutputChannel() {
			err := stream.Send(output)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// read messages from the stream and write the messages to corresponding channels
	// if the channel is not created, create the channel and invoke the reduceFn
	for {
		d, recvErr := stream.Recv()
		// if EOF, close all the channels
		if recvErr == io.EOF {
			taskManager.CloseAll()
			break
		}
		if recvErr != nil {
			// the error here is returned by stream.Recv()
			// it's already a gRPC error
			return recvErr
		}

		switch d.Operation.Event {
		case reducepb.ReduceRequest_WindowOperation_OPEN:
			// create a new reduce task and start the reduce operation
			err = taskManager.CreateTask(ctx, d, fs.Reducer)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		case reducepb.ReduceRequest_WindowOperation_CLOSE:
			// close the reduce task
			taskManager.CloseTask(d)
		case reducepb.ReduceRequest_WindowOperation_APPEND:
			// append the datum to the reduce task
			err = taskManager.AppendToTask(d, fs.Reducer)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		}
	}

	taskManager.WaitAll()
	// wait for the go routine which reads from the output channel and sends to the stream to return
	err = g.Wait()
	if err != nil {
		statusErr := status.Errorf(codes.Internal, err.Error())
		return statusErr
	}

	return nil
}
