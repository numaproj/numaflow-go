package function

import (
	"context"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct { // TODO: better name??...
	key       string
	value     []byte
	eventTime time.Time
	watermark time.Time
}

func (h *handlerDatum) Key() string {
	return h.key
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

// Service implements the proto gen server interface and contains the map operation handler and the reduce operation handler.
type Service struct {
	functionpb.UnimplementedUserDefinedFunctionServer

	Mapper  MapHandler
	Reducer ReduceHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*functionpb.ReadyResponse, error) {
	return &functionpb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a function to each datum element
func (fs *Service) MapFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	var hd = handlerDatum{
		key:       d.GetKey(),
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
	}
	messages, err := fs.Mapper.HandleDo(ctx, &hd)
	if err != nil {
		return nil, err
	}
	var elements []*functionpb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.Datum{
			Key:       m.Key,
			Value:     m.Value,
			EventTime: d.GetEventTime(),
			Watermark: d.GetWatermark(),
		})
	}
	datumList := &functionpb.DatumList{
		Elements: elements,
	}
	return datumList, nil
}

// ReduceFn applies a reduce function to a datum stream.
// TODO: implement ReduceFn
func (fs *Service) ReduceFn(fnServer functionpb.UserDefinedFunction_ReduceFnServer) error {
	ctx := fnServer.Context()
	_ = ctx
	return nil
}
