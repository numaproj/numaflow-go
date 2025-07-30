package reducer

import (
	"context"
	"errors"
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

	// Error group to manage goroutines, the groupCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the Wait returns.
	g, groupCtx := errgroup.WithContext(ctx)

	taskManager := newReduceTaskManager(ctx, fs.reducerCreatorHandle)

	// read from the response of the tasks and write to gRPC stream.
	g.Go(func() error {
		for {
			select {
			case <-groupCtx.Done():
				return nil
			case output, ok := <-taskManager.OutputChannel():
				if !ok {
					return nil
				}
				if err := stream.Send(output); err != nil {
					return status.Errorf(codes.Internal, "failed to send response: %v", err)
				}
			}
		}
	})

	// monitor error channel and handle task errors
	g.Go(func() error {
		for {
			select {
			case <-groupCtx.Done():
				return nil
			case errFromTask, ok := <-taskManager.ErrorChannel():
				if !ok {
					return nil
				}
				// As multiple streams will be running in parallel, we need to ensure that the shutdownCh is called only once.
				// Otherwise there could be a race condition where multiple streams try to call the shutdownCh.
				fs.once.Do(func() {
					log.Printf("Stopping the ReduceFn with err, %s", errFromTask)
					fs.shutdownCh <- struct{}{}
				})
				return errFromTask
			}
		}
	})

	var readErr error
readLoop:
	for {
		req, err := recvWithContext(groupCtx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the ReduceFn")
			break
		}
		if errors.Is(err, io.EOF) {
			log.Printf("EOF received, stopping the ReduceFn")
			taskManager.CloseAll()
			// wait for all tasks to complete and close output channel
			taskManager.WaitAll()
			// close error channel after all tasks are done
			taskManager.CloseErrorChannel()
			// cancel context to stop error group goroutines
			cancel()
			break
		}
		if err != nil {
			log.Printf("Failed to receive request: %v", err)
			// read loop is not part of the error group, so we need to cancel the context
			// to signal the other goroutines to stop processing.
			// this error is from inner stream wrapped in recvWithContext.
			cancel()
			readErr = err
			break
		}

		// for Aligned windows, its just open or append operation
		// close signal will be sent to all the reducers when grpc
		// input stream gets EOF.
		switch req.Operation.Event {
		case reducepb.ReduceRequest_WindowOperation_OPEN:
			// create a new reduce task and start the reduce operation
			err = taskManager.CreateTask(req)
			if err != nil {
				cancel()
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		case reducepb.ReduceRequest_WindowOperation_APPEND:
			// append the datum to the reduce task
			err = taskManager.AppendToTask(req)
			if err != nil {
				cancel()
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		}
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the ReduceFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
		})
		return err
	}

	if readErr != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the ReduceFn because of error while reading requests, %s", readErr)
			fs.shutdownCh <- struct{}{}
		})
		return readErr
	}
	return nil
}

// recvWithContext wraps stream.Recv() to respect context cancellation. We achieve that by writing to another channel and
// listening on the new channel also with ctx.Done so it can short-circuit if ctx is closed.
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
		return result.req, result.err
	}
}
