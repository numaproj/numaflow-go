package function

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
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
	var key string

	// get key from gPRC metadata
	if grpcMD, ok := metadata.FromIncomingContext(ctx); ok {
		keyValue := grpcMD.Get(DatumKey)
		if len(keyValue) > 1 {
			return nil, fmt.Errorf("expect extact one key but got %d keys", len(keyValue))
		} else if len(keyValue) == 1 {
			key = keyValue[0]
		} else {
			// do nothing: the length equals zero is valid, meaning the key is an empty string ""
		}
	}
	var hd = handlerDatum{
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
	}
	messages := fs.Mapper.HandleDo(ctx, key, &hd)
	var elements []*functionpb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.Datum{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	datumList := &functionpb.DatumList{
		Elements: elements,
	}
	return datumList, nil
}

// ReduceFn applies a reduce function to a datum stream.
func (fs *Service) ReduceFn(stream functionpb.UserDefinedFunction_ReduceFnServer) error {
	var (
		ctx      = stream.Context()
		key      string
		reduceCh = make(chan Datum)
		md       Metadata
	)

	// get key and metadata from gPRC metadata
	if grpcMD, ok := metadata.FromIncomingContext(ctx); ok {
		// get Key
		keyValue := grpcMD.Get(DatumKey)
		if len(keyValue) > 1 {
			return fmt.Errorf("expect extact one key but got %d keys", len(keyValue))
		} else if len(keyValue) == 1 {
			key = keyValue[0]
		} else {
			// do nothing: the length equals zero is valid, meaning the key is an empty string ""
		}
		// TODO: get metadata
	}

	var (
		datumList []*functionpb.Datum
		wg        sync.WaitGroup
	)

	go func() {
		wg.Add(1)
		defer wg.Done()
		messages := fs.Reducer.HandleDo(ctx, key, reduceCh, md)
		for _, msg := range messages {
			datumList = append(datumList, &functionpb.Datum{
				Key:   msg.Key,
				Value: msg.Value,
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
		if err != nil {
			return err
		}
		var hd = &handlerDatum{
			value:     datum.GetValue(),
			eventTime: datum.GetEventTime().EventTime.AsTime(),
			watermark: datum.GetWatermark().Watermark.AsTime(),
		}
		reduceCh <- hd
	}
}
