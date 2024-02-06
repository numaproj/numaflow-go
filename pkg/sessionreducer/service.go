package sessionreducer

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	sessionreducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/sessionreduce.sock"
	delimiter             = ":"
	serverInfoFilePath    = "/var/run/numaflow/sessionreducer-server-info"
)

// Service implements the proto gen server interface and contains the sesionreduce operation handler.
type Service struct {
	sessionreducepb.UnimplementedSessionReduceServer
	creatorHandle SessionReducerCreator
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sessionreducepb.ReadyResponse, error) {
	return &sessionreducepb.ReadyResponse{Ready: true}, nil
}

// SessionReduceFn applies a session reduce function to a request stream and streams the results.
func (fs *Service) SessionReduceFn(stream sessionreducepb.SessionReduce_SessionReduceFnServer) error {

	ctx := stream.Context()
	taskManager := newReduceTaskManager(fs.creatorHandle)
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

		// if the stream is closed, break and wait for the tasks to return
		if recvErr == io.EOF {
			break
		}

		if recvErr != nil {
			statusErr := status.Errorf(codes.Internal, recvErr.Error())
			return statusErr
		}

		// invoke the appropriate task manager method based on the operation
		switch d.Operation.Event {
		case sessionreducepb.SessionReduceRequest_WindowOperation_OPEN:
			// create a new task and start the session reduce operation
			// also append the datum to the task
			err := taskManager.CreateTask(ctx, d)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE:
			// close the task
			taskManager.CloseTask(d)
		case sessionreducepb.SessionReduceRequest_WindowOperation_APPEND:
			// append the datum to the task
			err := taskManager.AppendToTask(ctx, d)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_MERGE:
			// merge the tasks
			err := taskManager.MergeTasks(ctx, d)
			if err != nil {
				statusErr := status.Errorf(codes.Internal, err.Error())
				return statusErr
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND:
			// expand the task
			err := taskManager.ExpandTask(d)
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
