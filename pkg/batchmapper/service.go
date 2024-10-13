package batchmapper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/batchmap.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

// Service implements the proto gen server interface and contains the map operation handler.
type Service struct {
	mappb.UnimplementedMapServer
	BatchMapper BatchMapper
	shutdownCh  chan<- struct{}
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
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return fs.receiveRequests(ctx, stream, datumStreamCh)
		})

		g.Go(func() error {
			return fs.processData(ctx, stream, datumStreamCh)
		})

		// Wait for the goroutines to finish
		if err := g.Wait(); err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Stopping the BatchMapFn")
				return nil
			}
			log.Printf("Stopping the BatchMapFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
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

// receiveRequests receives the requests from the client and writes them to the datumStreamCh channel.
func (fs *Service) receiveRequests(ctx context.Context, stream mappb.Map_MapFnServer, datumStreamCh chan<- Datum) error {
	defer close(datumStreamCh)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("end of batch map stream")
			return err
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

		datumStreamCh <- datum
	}
	return nil
}

// processData invokes the batch mapper to process the data and sends the response back to the client.
func (fs *Service) processData(ctx context.Context, stream mappb.Map_MapFnServer, datumStreamCh chan Datum) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside batch map handler: %v %v", r, string(debug.Stack()))
			err = fmt.Errorf("panic inside batch map handler: %v", r)
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
			Status: &mappb.Status{
				Eot: true,
			},
		}
		if err := stream.Send(eot); err != nil {
			log.Println("BatchMapFn: Got an error while Send() on stream", err)
		}
	}
	return nil
}
