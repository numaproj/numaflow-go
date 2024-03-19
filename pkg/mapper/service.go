package mapper

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/map.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64
	serverInfoFilePath    = "/var/run/numaflow/mapper-server-info"
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	mappb.UnimplementedMapServer
	Mapper Mapper
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mappb.ReadyResponse, error) {
	return &mappb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(ctx context.Context, d *mappb.MapRequest) (*mappb.MapResponse, error) {
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().AsTime(), d.GetWatermark().AsTime(), d.GetHeaders())
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
