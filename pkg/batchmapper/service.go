package batchmapper

import (
	"context"
	"fmt"
	"io"
	"log"

	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	batchmappb "github.com/numaproj/numaflow-go/pkg/apis/proto/batchmap/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/batchmap.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	batchmappb.UnimplementedBatchMapServer
	BatchMapper BatchMapper
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*batchmappb.ReadyResponse, error) {
	return &batchmappb.ReadyResponse{Ready: true}, nil
}

// BatchMapFn applies a user defined function to a stream of request element and streams back the responses for them.
func (fs *Service) BatchMapFn(stream batchmappb.BatchMap_BatchMapFnServer) error {
	ctx := stream.Context()
	var g errgroup.Group

	// totalRequests is a counter for keeping a track of the number of datum requests
	// that were received on the stream. We use an atomic int as this needs to be synchronized
	// between the request/response go routines.
	totalRequests := atomic.NewInt32(0)

	// datumStreamCh is used to stream messages to the user code interface
	// As the BatchMap interface expects a list of request elements
	// we read all the requests coming on the stream and keep streaming them to the user code on this channel.
	datumStreamCh := make(chan Datum)

	// go routine to invoke the user handler function, and process the responses.
	g.Go(func() error {
		// Apply the user BatchMap implementation function
		responses := fs.BatchMapper.BatchMap(ctx, datumStreamCh)

		// If the number of responses received does not align with the request batch size,
		// we will not be able to process the data correctly.
		// This should be marked as an error and the container is restarted.
		// As this is a user error, we restart the container to mitigate any transient error otherwise, this
		// crash should indicate to the user that there is some issue.
		if len(responses.Items()) != int(totalRequests.Load()) {
			errMsg := fmt.Sprintf("batchMapFn: mismatch between length of batch requests and responses, "+
				"expected:%d, got:%d", int(totalRequests.Load()), len(responses.Items()))
			log.Panic(errMsg)
		}

		// iterate over the responses received and covert to the required proto format
		for _, batchResp := range responses.Items() {
			var elements []*batchmappb.BatchMapResponse_Result
			for _, resp := range batchResp.Items() {
				elements = append(elements, &batchmappb.BatchMapResponse_Result{
					Keys:  resp.Keys(),
					Value: resp.Value(),
					Tags:  resp.Tags(),
				})
			}
			singleRequestResp := &batchmappb.BatchMapResponse{
				Results: elements,
				Id:      batchResp.Id(),
			}
			// We stream back the result for a single request ID
			// this would contain all the responses for that request.
			err := stream.Send(singleRequestResp)
			if err != nil {
				log.Println("BatchMapFn: Got an error while Send() on stream", err)
				return err
			}
		}
		return nil
	})

	// loop to keep reading messages from the stream and sending it to the datumStreamCh
	for {
		d, err := stream.Recv()
		// if we see EOF on the stream we do not have any more messages coming up
		if err == io.EOF {
			// close the input data channel to indicate that no more messages expected
			close(datumStreamCh)
			break
		}
		if err != nil {
			// close the input data channel to indicate that no more messages expected
			close(datumStreamCh)
			log.Println("BatchMapFn: Got an error while recv() on stream", err)
			return err
		}
		var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders(), d.GetId(), d.GetKeys())
		// send the datum to the input channel
		datumStreamCh <- hd
		// Increase the counter for number of requests received
		totalRequests.Inc()
	}

	// wait for all the responses to be processed
	err := g.Wait()
	// if there was any error during processing return the error
	if err != nil {
		statusErr := status.Errorf(codes.Internal, err.Error())
		return statusErr
	}

	// Once all responses are sent we can return, this would indicate the end of the rpc and
	// send an EOF to the client on the stream
	return nil
}
