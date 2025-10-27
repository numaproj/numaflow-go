package sessionreducer

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
	shutdownCh    chan<- struct{}
	once          sync.Once
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sessionreducepb.ReadyResponse, error) {
	return &sessionreducepb.ReadyResponse{Ready: true}, nil
}

// SessionReduceFn applies a session reduce function to a request stream and streams the results.
func (fs *Service) SessionReduceFn(stream sessionreducepb.SessionReduce_SessionReduceFnServer) error {
	g, groupCtx := errgroup.WithContext(stream.Context())

	taskManager := newReduceTaskManager(groupCtx, fs.creatorHandle)

	// read from the response of the tasks and write to gRPC stream.
	g.Go(func() error {
		for {
			select {
			case <-groupCtx.Done():
				return nil
			case response, ok := <-taskManager.OutputChannel():
				if !ok {
					return nil
				}
				if err := stream.Send(response); err != nil {
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
				// Otherwise there could be a race condition where multiple streams try to send to the shutdownCh.
				fs.once.Do(func() {
					log.Printf("Stopping the SessionReduceFn with err, %s", errFromTask)
					select {
					case fs.shutdownCh <- struct{}{}:
						// signal enqueued
					default:
						log.Println("Shutdown signal already enqueued or watcher exited; skipping shutdown send")
					}
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
			log.Printf("Context cancelled, stopping the SessionReduceFn")
			break
		}
		if errors.Is(err, io.EOF) {
			log.Printf("EOF received, stopping the SessionReduceFn")
			taskManager.WaitAll()
			taskManager.CloseErrorChannel()
			break readLoop
		}
		if err != nil {
			log.Printf("Failed to receive request: %v", err)
			readErr = status.Errorf(codes.Internal, "failed to receive from stream: %v", err)
			break readLoop
		}

		// invoke the appropriate task manager method based on the operation
		switch req.Operation.Event {
		case sessionreducepb.SessionReduceRequest_WindowOperation_OPEN:
			// create a new task and start the session reduce operation
			// also append the datum to the task
			err := taskManager.CreateTask(groupCtx, req)
			if err != nil {
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE:
			// close the task
			taskManager.CloseTask(req)
		case sessionreducepb.SessionReduceRequest_WindowOperation_APPEND:
			// append the datum to the task
			err := taskManager.AppendToTask(groupCtx, req)
			if err != nil {
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_MERGE:
			// merge the tasks
			err := taskManager.MergeTasks(groupCtx, req)
			if err != nil {
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		case sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND:
			// expand the task
			err := taskManager.ExpandTask(req)
			if err != nil {
				readErr = status.Errorf(codes.Internal, "%s", err.Error())
				break readLoop
			}
		}
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the SessionReduceFn with err, %s", err)
			select {
			case fs.shutdownCh <- struct{}{}:
				// signal enqueued
			default:
				log.Println("Shutdown signal already enqueued or watcher exited; skipping shutdown send")
			}
		})
		return err
	}

	if readErr != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the SessionReduceFn because of error while reading requests, %s", readErr)
			select {
			case fs.shutdownCh <- struct{}{}:
				// signal enqueued
			default:
				log.Println("Shutdown signal already enqueued or watcher exited; skipping shutdown send")
			}
		})
		return readErr
	}
	return nil
}

// recvWithContext wraps stream.Recv() with context cancellation support
func recvWithContext(ctx context.Context, stream sessionreducepb.SessionReduce_SessionReduceFnServer) (*sessionreducepb.SessionReduceRequest, error) {
	type result struct {
		req *sessionreducepb.SessionReduceRequest
		err error
	}

	resultCh := make(chan result, 1)
	go func() {
		req, err := stream.Recv()
		resultCh <- result{req: req, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultCh:
		return res.req, res.err
	}
}
