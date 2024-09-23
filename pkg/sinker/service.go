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
	req, err := stream.Recv()
	if err != nil {
		log.Printf("error receiving handshake from stream: %v", err)
		return err
	}

	if req.Handshake == nil || !req.Handshake.Sot {
		return fmt.Errorf("expected handshake message")
	}

	// Send handshake response
	handshakeResponse := &sinkpb.SinkResponse{
		Result: &sinkpb.SinkResponse_Result{
			Status: sinkpb.Status_SUCCESS,
		},
		Handshake: &sinkpb.Handshake{
			Sot: true,
		},
	}
	if err := stream.Send(handshakeResponse); err != nil {
		return err
	}

	for {
		datumStreamCh := make(chan Datum)
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			defer close(datumStreamCh)
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic inside sink handler: %v %v", r, string(debug.Stack()))
					fs.shutdownCh <- struct{}{}
				}
			}()

			for {
				// Receive sink requests from the stream
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
					// End of transmission, break to start a new sink invocation
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

				// Send datum to the channel
				datumStreamCh <- datum
			}
			return nil
		})

		// invoke the sink function, and send the responses back to the client
		g.Go(func() error {
			responses := fs.Sinker.Sink(ctx, datumStreamCh)
			for _, response := range responses {
				var status sinkpb.Status
				if response.Fallback {
					status = sinkpb.Status_FALLBACK
				} else if response.Success {
					status = sinkpb.Status_SUCCESS
				} else {
					status = sinkpb.Status_FAILURE
				}

				sinkResponse := &sinkpb.SinkResponse{
					Result: &sinkpb.SinkResponse_Result{
						Id:     response.ID,
						Status: status,
						ErrMsg: response.Err,
					},
				}
				if err := stream.Send(sinkResponse); err != nil {
					log.Printf("error sending sink response: %v", err)
					return err
				}
			}
			return nil
		})

		// Wait for the goroutines to finish
		err := g.Wait()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
