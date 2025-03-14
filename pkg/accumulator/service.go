package accumulator

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

	accumulatorpb "github.com/numaproj/numaflow-go/pkg/apis/proto/accumulator/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/accumulator.sock"
	serverInfoFilePath    = "/var/run/numaflow/accumulator-server-info"
)

var errAccumulatorPanic = errors.New("UDF_EXECUTION_ERROR(accumulator)")

// Service implements the proto gen server interface and contains the accumulator operation handler.
type Service struct {
	accumulatorpb.UnimplementedAccumulatorServer
	accumulator AccumulatorCreator
	once        sync.Once
	shutdownCh  chan<- struct{}
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

	taskManager := newAccumulateTaskManager(groupCtx, g, fs.accumulator)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
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
			cancel()
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
	return nil
}

// recvWithContext wraps stream.Recv() to respect context cancellation.
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
