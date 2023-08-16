package _map

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	v1.UnimplementedMapServer
	Mapper Mapper
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(ctx context.Context, d *v1.MapRequest) (*v1.MapResponseList, error) {
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime())
	messages := fs.Mapper.Map(ctx, d.GetKeys(), hd)
	var elements []*v1.MapResponse
	for _, m := range messages.Items() {
		elements = append(elements, &v1.MapResponse{
			Keys:  m.Keys(),
			Value: m.Value(),
			Tags:  m.Tags(),
		})
	}
	datumList := &v1.MapResponseList{
		Elements: elements,
	}
	return datumList, nil
}
