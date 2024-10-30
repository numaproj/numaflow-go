package sinker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
)

const (
	uds                     = "unix"
	defaultMaxMessageSize   = 1024 * 1024 * 64 // 64MB
	address                 = "/var/run/numaflow/sink.sock"
	fbAddress               = "/var/run/numaflow/fb-sink.sock"
	serverInfoFilePath      = "/var/run/numaflow/sinker-server-info"
	fbServerInfoFilePath    = "/var/run/numaflow/fb-sinker-server-info"
	EnvUDContainerType      = "NUMAFLOW_UD_CONTAINER_TYPE"
	UDContainerFallbackSink = "fb-udsink"
)

// handlerDatum implements the Datum interface and is used in the sink functions.
type handlerDatum struct {
	id        string
	keys      []string
	value     []byte
	eventTime time.Time
	watermark time.Time
	headers   map[string]string
}

func (h *handlerDatum) Keys() []string {
	return h.keys
}

func (h *handlerDatum) ID() string {
	return h.id
}

func (h *handlerDatum) Value() []byte {
	return h.value
}

func (h *handlerDatum) EventTime() time.Time {
	return h.eventTime
}

func (h *handlerDatum) Watermark() time.Time {
	return h.watermark
}

func (h *handlerDatum) Headers() map[string]string {
	return h.headers
}

// Service implements the proto gen server interface and contains the sinkfn operation handler.
type Service struct {
	sinkpb.UnimplementedSinkServer
	shutdownCh chan<- struct{}
	Sinker     Sinker
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sinkpb.ReadyResponse, error) {
	return &sinkpb.ReadyResponse{Ready: true}, nil
}

// SinkFn applies a sink function to a every element.
func (fs *Service) SinkFn(stream sinkpb.Sink_SinkFnServer) error {
	ctx := stream.Context()

	// Perform handshake before entering the main loop
	if err := fs.performHandshake(stream); err != nil {
		return err
	}

	for {
		datumStreamCh := make(chan Datum)
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return fs.receiveRequests(stream, datumStreamCh)
		})

		g.Go(func() error {
			return fs.processData(ctx, stream, datumStreamCh)
		})

		// Wait for the goroutines to finish
		if err := g.Wait(); err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Stopping the SinkFn")
				return nil
			}
			log.Printf("Stopping the SinkFn with err, %s", err)
			fs.shutdownCh <- struct{}{}
			return err
		}
	}
}

// performHandshake performs the handshake with the client.
func (fs *Service) performHandshake(stream sinkpb.Sink_SinkFnServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("error receiving handshake from stream: %v", err)
		return err
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("expected handshake message")
	}

	handshakeResponse := &sinkpb.SinkResponse{
		Handshake: &sinkpb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return err
	}

	return nil
}

// receiveRequests receives the requests from the client writes them to the datumStreamCh channel.
func (fs *Service) receiveRequests(stream sinkpb.Sink_SinkFnServer, datumStreamCh chan<- Datum) error {
	defer close(datumStreamCh)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("end of sink stream")
			return err
		}
		if err != nil {
			log.Printf("error receiving from sink stream: %v", err)
			return err
		}

		if req.Status != nil && req.Status.Eot {
			break
		}

		datum := &handlerDatum{
			id:        req.GetRequest().GetId(),
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

// processData invokes the sinker to process the data and sends the response back to the client.
func (fs *Service) processData(ctx context.Context, stream sinkpb.Sink_SinkFnServer, datumStreamCh chan Datum) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside sink handler: %v %v", r, string(debug.Stack()))
			err = fmt.Errorf("panic inside sink handler: %v", r)
		}
	}()
	responses := fs.Sinker.Sink(ctx, datumStreamCh)

	var resultList []*sinkpb.SinkResponse_Result
	for _, msg := range responses {
		if msg.Fallback {
			resultList = append(resultList, &sinkpb.SinkResponse_Result{
				Id:     msg.ID,
				Status: sinkpb.Status_FALLBACK,
			})
		} else if msg.Success {
			resultList = append(resultList, &sinkpb.SinkResponse_Result{
				Id:     msg.ID,
				Status: sinkpb.Status_SUCCESS,
			})
		} else {
			resultList = append(resultList, &sinkpb.SinkResponse_Result{
				Id:     msg.ID,
				Status: sinkpb.Status_FAILURE,
				ErrMsg: msg.Err,
			})
		}
	}
	if err := stream.Send(&sinkpb.SinkResponse{
		Results: resultList,
	}); err != nil {
		log.Printf("error sending sink response: %v", err)
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// send the end of transmission message
	eotResponse := &sinkpb.SinkResponse{
		Status: &sinkpb.TransmissionStatus{Eot: true},
	}
	if err := stream.Send(eotResponse); err != nil {
		log.Printf("error sending end of transmission message: %v", err)
		return err
	}
	return nil
}
