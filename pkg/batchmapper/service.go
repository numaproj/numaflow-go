package batchmapper

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

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/batchmap.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

var errBatchMapHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// Service implements the proto gen server interface and contains the map operation handler.
type Service struct {
	mappb.UnimplementedMapServer
	BatchMapper BatchMapper
	shutdownCh  chan<- struct{}
	once        sync.Once
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to a stream of request elements and streams back the responses for them.
func (fs *Service) MapFn(stream mappb.Map_MapFnServer) error {
	ctx := stream.Context()

	// Perform handshake before entering the main loop
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	for {
		datumStreamCh := make(chan Datum)
		g, groupCtx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return fs.receiveRequests(groupCtx, stream, datumStreamCh)
		})

		g.Go(func() error {
			return fs.processData(groupCtx, stream, datumStreamCh)
		})

		// Wait for the goroutines to finish
		if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
			fs.once.Do(func() {
				log.Printf("Stopping the BatchMapFn with err, %s", err)
				fs.shutdownCh <- struct{}{}
			})
			return err
		}
	}
}

// performHandshake performs the handshake with the client.
func (fs *Service) performHandshake(stream mappb.Map_MapFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("error receiving handshake from stream: %v", err)
		return err
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("expected handshake message")
	}

	handshakeResponse := &mappb.MapResponse{
		Handshake: &mappb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return err
	}

	return nil
}

// recvWithContext wraps stream.Recv() to respect context cancellation.
func recvWithContext(ctx context.Context, stream mappb.Map_MapFnServer) (*mappb.MapRequest, error) {
	type recvResult struct {
		req *mappb.MapRequest
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

// receiveRequests receives the requests from the client and writes them to the datumStreamCh channel.
func (fs *Service) receiveRequests(ctx context.Context, stream mappb.Map_MapFnServer, datumStreamCh chan<- Datum) error {
	defer close(datumStreamCh)

	for {
		req, err := recvWithContext(ctx, stream)
		if errors.Is(err, context.Canceled) {
			log.Printf("Context cancelled, stopping the MapBatchFn")
			return err
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			log.Printf("error receiving from batch map stream: %v", err)
			return err
		}

		if req.Status != nil && req.Status.Eot {
			break
		}

		datum := &handlerDatum{
			id:        req.GetId(),
			value:     req.GetRequest().GetValue(),
			keys:      req.GetRequest().GetKeys(),
			eventTime: req.GetRequest().GetEventTime().AsTime(),
			watermark: req.GetRequest().GetWatermark().AsTime(),
			headers:   req.GetRequest().GetHeaders(),
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case datumStreamCh <- datum:
		}
	}
	return nil
}

// processData invokes the batch mapper to process the data and sends the response back to the client.
func (fs *Service) processData(ctx context.Context, stream mappb.Map_MapFnServer, datumStreamCh chan Datum) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside batch map handler: %v %v", r, string(debug.Stack()))
			st, _ := status.Newf(codes.Internal, "%s: %v", errBatchMapHandlerPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
		}
	}()

	responses := fs.BatchMapper.BatchMap(ctx, datumStreamCh)

	for _, batchResp := range responses.Items() {
		var elements []*mappb.MapResponse_Result
		for _, resp := range batchResp.Items() {
			elements = append(elements, &mappb.MapResponse_Result{
				Keys:  resp.Keys(),
				Value: resp.Value(),
				Tags:  resp.Tags(),
			})
		}
		singleRequestResp := &mappb.MapResponse{
			Results: elements,
			Id:      batchResp.Id(),
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := stream.Send(singleRequestResp); err != nil {
			log.Println("BatchMapFn: Got an error while Send() on stream", err)
			return err
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// send the end of transmission message
		eot := &mappb.MapResponse{
			Status: &mappb.TransmissionStatus{
				Eot: true,
			},
		}
		if err := stream.Send(eot); err != nil {
			log.Println("BatchMapFn: Got an error while Send() on stream", err)
		}
	}
	return nil
}
