package function

import (
	"context"
	"io"
	"sync"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (fs *Service) ReduceFn(stream functionpb.UserDefinedFunction_ReduceFnServer) error {
	var ctx = context.Background() // TODO: revisit ctx
	var reduceCh = make(chan Datum)
	var md Metadata

	var datumList []*functionpb.Datum
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		messages, err := fs.Reducer.HandleDo(ctx, reduceCh, md)
		if err != nil {
			// will return an empty datumList
			return
		}
		for _, msg := range messages {
			datumList = append(datumList, &functionpb.Datum{
				Key:       msg.Key,
				Value:     msg.Value,
				EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})}, // TODO: what's the correct value?...
				Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
			})
		}
	}()

	for {
		datum, err := stream.Recv()
		if err == io.EOF {
			close(reduceCh)
			wg.Wait()
			return stream.SendAndClose(&functionpb.DatumList{
				Elements: datumList,
			})
		}

		var hd = &handlerDatum{
			key:       datum.GetKey(),
			value:     datum.GetValue(),
			eventTime: datum.GetEventTime().EventTime.AsTime(),
			watermark: datum.GetWatermark().Watermark.AsTime(),
		}
		reduceCh <- hd
	}
}
