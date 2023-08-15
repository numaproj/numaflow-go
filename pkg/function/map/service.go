package mapsvc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/mapfn"
	"github.com/numaproj/numaflow-go/pkg/function"
)

// Service implements the proto gen server interface and contains the map operation
// handler.
type Service struct {
	mapfn.UnimplementedMapServer
	Mapper function.MapHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*mapfn.ReadyResponse, error) {
	return &mapfn.ReadyResponse{Ready: true}, nil
}

// MapFn applies a user defined function to each request element and returns a list of results.
func (fs *Service) MapFn(ctx context.Context, d *mapfn.MapRequest) (*mapfn.MapResponseList, error) {
	var hd = function.NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime())
	messages := fs.Mapper.HandleDo(ctx, d.GetKeys(), hd)
	var elements []*mapfn.MapResponse
	for _, m := range messages.Items() {
		elements = append(elements, &mapfn.MapResponse{
			Keys:  m.Keys(),
			Value: m.Value(),
			Tags:  m.Tags(),
		})
	}
	datumList := &mapfn.MapResponseList{
		Elements: elements,
	}
	return datumList, nil
}
