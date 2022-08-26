package function

import (
	"context"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

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

// DoFn applies a function to each datum element
func (fs *Service) DoFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	messages, err := fs.Mapper.HandleDo(ctx, d.GetKey(), d.GetValue())
	if err != nil {
		return nil, err
	}
	var elements []*functionpb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.Datum{
			Key:            m.Key,
			Value:          m.Value,
			EventTime:      d.GetEventTime(),
			PaneInfo:       d.GetPaneInfo(),
			IntervalWindow: d.GetIntervalWindow(),
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
	return nil
}
