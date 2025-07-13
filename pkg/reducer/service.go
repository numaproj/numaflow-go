package reducer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

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
	serverInfoFilePath    = "/var/run/numaflow/reducer-server-info"
)

// Service implements the proto gen server interface and contains the reduce operation handler.
type Service struct {
	reducepb.UnimplementedReduceServer
	reducerCreatorHandle ReducerCreator
	shutdownCh           chan<- struct{}
	once                 sync.Once
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*reducepb.ReadyResponse, error) {
	return &reducepb.ReadyResponse{Ready: true}, nil
}

// ReduceFn applies a reduce function to a request stream and returns a list of results.
func (fs *Service) ReduceFn(stream reducepb.Reduce_ReduceFnServer) error {

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	g, groupCtx := errgroup.WithContext(ctx)
	taskManager := newReduceTaskManager(fs.reducerCreatorHandle)

	// Goroutine to send output to stream or return error if any.
	g.Go(func() error {
		for {
			select {
			case output := <-taskManager.OutputChannel():
				sendErr := stream.Send(output)
				if sendErr != nil {
					return sendErr
				}
			case err := <-taskManager.ErrorChannel():
				if err != nil {
					log.Printf("received err in ReduceFn: %v", err)
				}
				return err
			case <-groupCtx.Done():
				return groupCtx.Err()
			}
		}
	})

	var readErr error
outer:
	for {
		req, err := recvWithContext(groupCtx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the ReduceFn")
			break outer
		}
		if errors.Is(err, io.EOF) {
			taskManager.CloseAll()
			break outer
		}
		if err != nil {
			readErr = err
			cancel()
			break outer
		}

		// Handle the request based on the operation type
		g.Go(func() error {
			switch req.Operation.Event {
			case reducepb.ReduceRequest_WindowOperation_OPEN:
				// create a new reduce task and start the reduce operation
				err := taskManager.CreateTask(groupCtx, req)
				if err != nil {
					return err
				}
			case reducepb.ReduceRequest_WindowOperation_APPEND:
				// append the datum to the reduce task
				err := taskManager.AppendToTask(groupCtx, req)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	// wait for all tasks to complete
	taskManager.WaitAll(groupCtx)

	// wait for the go routine which reads from the output channel and sends to the stream to return
	if err := g.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the ReduceFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
		})
		return err
	}
	if readErr != nil {
		return status.Errorf(codes.Internal, "%s", readErr.Error())
	}

	return nil
}

func recvWithContext(ctx context.Context, stream reducepb.Reduce_ReduceFnServer) (*reducepb.ReduceRequest, error) {
	type recvResult struct {
		req *reducepb.ReduceRequest
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
		fmt.Println("received request:", result.req)
		return result.req, result.err
	}
}
