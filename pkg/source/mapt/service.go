package mapt

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/source/transformerfn"
	"github.com/numaproj/numaflow-go/pkg/source"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service implements the proto gen server interface and contains the transformer operation
// handler.
type Service struct {
	transformerfn.UnimplementedSourceTransformerServer
	MapperT source.MapTHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*transformerfn.ReadyResponse, error) {
	return &transformerfn.ReadyResponse{Ready: true}, nil
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
// MapTFn can be used only at source vertex by source data transformer.
func (fs *Service) MapTFn(ctx context.Context, d *transformerfn.SourceTransformerRequest) (*transformerfn.SourceTransformerResponseList, error) {
	var hd = source.NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime())
	messageTs := fs.MapperT.HandleDo(ctx, d.GetKeys(), hd)
	var elements []*transformerfn.SourceTransformerResponse
	for _, m := range messageTs.Items() {
		elements = append(elements, &transformerfn.SourceTransformerResponse{
			EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(m.EventTime())},
			Keys:      m.Keys(),
			Value:     m.Value(),
			Tags:      m.Tags(),
		})
	}
	responseList := &transformerfn.SourceTransformerResponseList{
		Elements: elements,
	}
	return responseList, nil
}
