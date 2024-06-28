package mapper

import (
	"context"
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
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders(), emptyId)
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

func (fs *Service) BatchMapFn(stream mappb.Map_BatchMapFnServer) error {
	ctx := stream.Context()
	datums := make([]Datum, 0)
	for {
		d, err := stream.Recv()
		if err == io.EOF {
			log.Println("MYDEBUG: Got EOF here", err)
			break
		}
		if err != nil {
			log.Println("MYDEBUG: Got error here", err)
			return err
		}
		var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders(), d.GetId())
		datums = append(datums, hd)
	}

	index := 0
	responses := fs.BatchMapper.BatchMap(ctx, datums)
	for _, batchResp := range responses.Items() {
		var elements []*mappb.BatchMapResponse_Result
		for _, resp := range batchResp.Items() {
			elements = append(elements, &mappb.BatchMapResponse_Result{
				Keys:  resp.Keys(),
				Value: resp.Value(),
				Tags:  resp.Tags(),
			})
			index += 1
		}
		datumList := &mappb.BatchMapResponse{
			Results: elements,
			Id:      batchResp.Id(),
		}
		err := stream.Send(datumList)
		if err != nil {
			log.Println("MYDEBUG: Got error here - 2", err)
			return err
		}
	}
	log.Println("MYDEBUG: SENDING BACK FROM UDF ", index)
	return nil
}
