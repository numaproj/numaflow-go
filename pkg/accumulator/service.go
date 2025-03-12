package accumulator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"

	"golang.org/x/sync/errgroup"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	accumulator Accumulator
	once        sync.Once
	shutdownCh  chan<- struct{}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*accumulatorpb.ReadyResponse, error) {
	return &accumulatorpb.ReadyResponse{Ready: true}, nil
}

func (fs *Service) AccumulateFn(stream accumulatorpb.Accumulator_AccumulateFnServer) error {
	ctx := stream.Context()
	// perform handshake with client before processing requests
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	// Use error group to manage goroutines, the groupCtx is cancelled when any of the
	// goroutines return an error for the first time or the first time the wait returns.
	g, groupCtx := errgroup.WithContext(ctx)

	// Channel to collect responses
	responseCh := make(chan *accumulatorpb.AccumulatorResponse, 500) // FIXME: identify the right buffer size
	defer close(responseCh)

	datumInputCh := make(chan Datum, 500)
	datumOutputCh := make(chan Datum, 500)

	// Start the accumulator
	g.Go(func() error {
		if err := invokeAccumulator(groupCtx, fs.accumulator, datumInputCh, datumOutputCh); err != nil {
			log.Printf("Accumulator handler failed: %v", err)
			return err
		}
		close(datumOutputCh)
		return nil
	})

	// Start the accumulator response sender
	g.Go(func() error {
		for {
			select {
			case <-groupCtx.Done():
				return nil
			case datum, ok := <-datumOutputCh:
				if !ok {
					return nil
				}
				responseCh <- &accumulatorpb.AccumulatorResponse{
					Payload: &accumulatorpb.Payload{
						Keys:      datum.Keys(),
						Value:     datum.Value(),
						EventTime: timestamppb.New(datum.EventTime()),
						Watermark: timestamppb.New(datum.Watermark()),
						Headers:   datum.Headers(),
					},
				}
			}
		}
	})

	// Start the accumulator request sender
	g.Go(func() error {
		for {
			req, err := recvWithContext(groupCtx, stream)
			if errors.Is(err, context.Canceled) {
				log.Printf("Context cancelled, stopping the AccumulateFn")
				return nil
			}
			if errors.Is(err, io.EOF) {
				log.Printf("EOF received, stopping the AccumulateFn")
				return nil
			}
			if err != nil {
				log.Printf("Failed to receive request: %v", err)
				// read loop is not part of the error group, so we need to cancel the context
				// to signal the other goroutines to stop processing.
				cancel()
				return err
			}
			datumInputCh <- &handlerDatum{
				value:     req.GetPayload().GetValue(),
				eventTime: req.GetPayload().GetEventTime().AsTime(),
				watermark: req.GetPayload().GetWatermark().AsTime(),
				keys:      req.GetPayload().GetKeys(),
				headers:   req.GetPayload().GetHeaders(),
			}
		}
	})

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

// performHandshake handles the handshake logic at the start of the stream.
func (fs *Service) performHandshake(stream accumulatorpb.Accumulator_AccumulateFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive handshake: %v", err)
	}
	if req.GetHandshake() == nil || !req.GetHandshake().GetSot() {
		return status.Errorf(codes.InvalidArgument, "invalid handshake")
	}
	handshakeResponse := &accumulatorpb.AccumulatorResponse{
		Handshake: &accumulatorpb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return fmt.Errorf("sending handshake response to client over gRPC stream: %w", err)
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

// invokeAccumulator invokes the accumulator handler and recovers from panics.
func invokeAccumulator(ctx context.Context, accumulator Accumulator, inputDatumCh chan Datum, outputDatumCh chan Datum) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside accumulator handler: %v %v", r, string(debug.Stack()))
			st, _ := status.Newf(codes.Internal, "%s: %v", errAccumulatorPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
		}
	}()

	// Start the accumulator
	accumulator.Accumulate(ctx, inputDatumCh, outputDatumCh)
	return
}
