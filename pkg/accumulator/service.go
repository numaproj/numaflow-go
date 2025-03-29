package accumulator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	accumulatorpb "github.com/numaproj/numaflow-go/pkg/apis/proto/accumulator/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/accumulator.sock"
	serverInfoFilePath    = "/var/run/numaflow/accumulator-server-info"
)

var containerType = func() string {
	if val, exists := os.LookupEnv(shared.EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

var errAccumulatorPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", containerType)

// Service implements the proto gen server interface and contains the accumulator operation handler.
type Service struct {
	accumulatorpb.UnimplementedAccumulatorServer
	AccumulatorCreator AccumulatorCreator
	once               sync.Once
	shutdownCh         chan<- struct{}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*accumulatorpb.ReadyResponse, error) {
	return &accumulatorpb.ReadyResponse{Ready: true}, nil
}

func (fs *Service) AccumulateFn(stream accumulatorpb.Accumulator_AccumulateFnServer) error {
	ctx := stream.Context()

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Use error group to manage goroutines, the groupCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the wait returns.
	g, groupCtx := errgroup.WithContext(ctx)

	taskManager := newAccumulatorTaskManager(groupCtx, g, fs.AccumulatorCreator)

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

	var readErr error
	for {
		req, err := recvWithContext(groupCtx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the AccumulateFn")
			break
		}
		if errors.Is(err, io.EOF) {
			log.Printf("EOF received, stopping the AccumulateFn")
			taskManager.CloseAll()
			return nil
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

		switch req.Operation.Event {
		case accumulatorpb.AccumulatorRequest_WindowOperation_OPEN:
			taskManager.CreateTask(req)
		case accumulatorpb.AccumulatorRequest_WindowOperation_CLOSE:
			taskManager.CloseTask(req)
		case accumulatorpb.AccumulatorRequest_WindowOperation_APPEND:
			taskManager.AppendToTask(req)
		}
	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the AccumulateFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
		})
		return err
	}

	if readErr != nil {
		fs.once.Do(func() {
			log.Printf("Stopping the AccumulateFn because of error while reading requests, %s", readErr)
			fs.shutdownCh <- struct{}{}
		})
		return readErr
	}
	return nil
}

// recvWithContext wraps stream.Recv() to respect context cancellation. We achieve that by writing to another channel and
// listening on the new channel also with ctx.Done so it can short-circuit if ctx is closed. Without this approach, it would
// be blocking because nobody is listening on ctx.Done().
func recvWithContext(ctx context.Context, stream accumulatorpb.Accumulator_AccumulateFnServer) (*accumulatorpb.AccumulatorRequest, error) {
	type recvResult struct {
		req *accumulatorpb.AccumulatorRequest
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
