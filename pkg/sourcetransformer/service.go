package sourcetransformer

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransformer/v1"
)

// Service implements the proto gen server interface and contains the transformer operation
// handler.
type Service struct {
	v1.UnimplementedSourceTransformerServer
	Transformer SourceTransformer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*v1.ReadyResponse, error) {
	return &v1.ReadyResponse{Ready: true}, nil
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
// MapTFn can be used only at source vertex by source data transformer.
func (fs *Service) MapTFn(ctx context.Context, d *v1.SourceTransformerRequest) (*v1.SourceTransformerResponseList, error) {
	var hd = NewHandlerDatum(d.GetValue(), d.GetEventTime().EventTime.AsTime(), d.GetWatermark().Watermark.AsTime())
	messageTs := fs.Transformer.SourceTransformer(ctx, d.GetKeys(), hd)
	var elements []*v1.SourceTransformerResponse
	for _, m := range messageTs.Items() {
		elements = append(elements, &v1.SourceTransformerResponse{
			EventTime: &v1.EventTime{EventTime: timestamppb.New(m.EventTime())},
			Keys:      m.Keys(),
			Value:     m.Value(),
			Tags:      m.Tags(),
		})
	}
	responseList := &v1.SourceTransformerResponseList{
		Elements: elements,
	}
	return responseList, nil
}
