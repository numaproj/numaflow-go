package globalreducer

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	globalreducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/globalreduce/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/globalreduce.sock"
	delimiter             = ":"
)

// Service implements the proto gen server interface and contains the globalreduce operation handler.
type Service struct {
	globalreducepb.UnimplementedGlobalReduceServer
	globalReducer GlobalReducer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*globalreducepb.ReadyResponse, error) {
	return &globalreducepb.ReadyResponse{Ready: true}, nil
}

// GlobalReduceFn applies a global reduce function to a request stream and streams the results.
func (fs *Service) GlobalReduceFn(stream globalreducepb.GlobalReduce_GlobalReduceFnServer) error {

	ctx := stream.Context()
	taskManager := newReduceTaskManager(fs.globalReducer)
	// err group for the go routine which reads from the output channel and sends to the stream
	var g errgroup.Group

	g.Go(func() error {
		for output := range taskManager.OutputChannel() {
			err := stream.Send(output)
			if err != nil {
				return err
			}
		}
		return nil
	})

	for {
		d, recvErr := stream.Recv()

		// if the stream is closed, close all the tasks and break
		if recvErr == io.EOF {
			taskManager.CloseAll()
			break
		}

		if recvErr != nil {
			statusErr := status.Errorf(codes.Internal, recvErr.Error())
			return statusErr
		}

		// invoke the appropriate task manager method based on the operation
		switch d.Operation.Event {
		case globalreducepb.GlobalReduceRequest_WindowOperation_OPEN:
			// create a new task and start the global reduce operation
			// also append the datum to the task
			err := taskManager.CreateTask(ctx, d)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		case globalreducepb.GlobalReduceRequest_WindowOperation_CLOSE:
			// close the task
			taskManager.CloseTask(d)
		case globalreducepb.GlobalReduceRequest_WindowOperation_APPEND:
			// append the datum to the task
			err := taskManager.AppendToTask(ctx, d)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		}

	}

	// wait for all the tasks to return
	taskManager.WaitAll()

	// wait for the go routine which reads from the output channel and sends to the stream to return
	err := g.Wait()
	if err != nil {
		statusErr := status.Errorf(codes.Internal, err.Error())
		return statusErr
	}

	return nil
}
