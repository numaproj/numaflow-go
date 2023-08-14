package mapsvc

import (
	"context"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/map/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service implements the proto gen server interface and contains the map operation
// handler and the reduce operation handler.
type Service struct {
	v1.UnimplementedMapServer
	Mapper function.MapHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

// MapFn applies a function to each datum element.
func (fs *Service) MapFn(ctx context.Context, d *v1.MapRequest) (*v1.MapResponseList, error) {
	var hd = function.NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime(), nil)
	messages := fs.Mapper.HandleDo(ctx, d.GetKeys(), hd)
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
