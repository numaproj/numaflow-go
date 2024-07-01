package mapper

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/map.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
	emptyId               = ""
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	mappb.UnimplementedMapServer
	Mapper      Mapper
	BatchMapper BatchMapper
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(ctx context.Context, d *mappb.MapRequest) (*mappb.MapResponse, error) {
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders(), emptyId, d.GetKeys())
	messages := fs.Mapper.Map(ctx, d.GetKeys(), hd)
	var elements []*mappb.MapResponse_Result
	for _, m := range messages.Items() {
		elements = append(elements, &mappb.MapResponse_Result{
			Keys:  m.Keys(),
			Value: m.Value(),
			Tags:  m.Tags(),
		})
	}
	datumList := &mappb.MapResponse{
		Results: elements,
	}
	return datumList, nil
}

// MapStreamFn applies a user defined function to a stream of request element and streams back the responses for them.
func (fs *Service) MapStreamFn(stream mappb.Map_MapStreamFnServer) error {
	ctx := stream.Context()

	// As the BatchMap interface expects a list of request elements
	// we read all the requests coming on the stream and keep appending them together
	// and then finally send the array for processing.
	datums := make([]Datum, 0)
	for {
		d, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("MapStreamFn: Got an error while recv() on stream", err)
			return err
		}
		var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders(), d.GetId(), d.GetKeys())
		datums = append(datums, hd)
	}

	// Apply the user BatchMap implementation function
	responses := fs.BatchMapper.BatchMap(ctx, datums)

	// If the number of responses received does not align with the request batch size,
	// we will not be able to process the data correctly.
	// This should be marked as an error and shown to the user.
	// TODO(map-batch): We could potentially panic here as well
	if len(responses.Items()) != len(datums) {
		errMsg := "batchMapFn: mismatch between length of batch requests and responses"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	// iterate over the responses received and covert to the required proto format
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
		// We stream back the result for a single request ID
		// this would contain all the responses for that request.
		err := stream.Send(singleRequestResp)
		if err != nil {
			log.Println("MapStreamFn: Got an error while Send() on stream", err)
			return err
		}
	}
	// Once all responses are sent we can return, this would indicate the end of the rpc and
	// send an EOF to the client on the stream
	return nil
}
