package function

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// intervalWindow implements IntervalWindow interface which will be passed as metadata
// to reduce handlers
type intervalWindow struct {
	startTime time.Time
	endTime   time.Time
}

func NewIntervalWindow(startTime time.Time, endTime time.Time) IntervalWindow {
	return &intervalWindow{
		startTime: startTime,
		endTime:   endTime,
	}
}

func (i *intervalWindow) StartTime() time.Time {
	return i.startTime
}

func (i *intervalWindow) EndTime() time.Time {
	return i.endTime
}

// metadata implements Metadata interface which will be passed to reduce handlers
type metadata struct {
	intervalWindow IntervalWindow
}

func NewMetadata(window IntervalWindow) Metadata {
	return &metadata{intervalWindow: window}
}

func (m *metadata) IntervalWindow() IntervalWindow {
	return m.intervalWindow
}

// Service implements the proto gen server interface and contains the map operation handler and the reduce operation handler.
type Service struct {
	functionpb.UnimplementedUserDefinedFunctionServer

	Mapper  MapHandler
	MapperT MapTHandler
	Reducer ReduceHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*functionpb.ReadyResponse, error) {
	return &functionpb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a function to each datum element without modifying event time.
func (fs *Service) MapFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	var key string

	// get key from gPRC metadata
	if grpcMD, ok := grpcmd.FromIncomingContext(ctx); ok {
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

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
func (fs *Service) MapTFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	var key string

	// get key from gPRC metadata
	if grpcMD, ok := grpcmd.FromIncomingContext(ctx); ok {
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
	messages := fs.MapperT.HandleDo(ctx, key, &hd)
	var elements []*functionpb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.Datum{
			EventTime: &functionpb.EventTime{EventTime: timestamppb.New(m.EventTime)},
			Key:       m.Key,
			Value:     m.Value,
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
		key       string
		md        Metadata
		err       error
		startTime int64
		endTime   int64
		reduceCh  = make(chan Datum)
		ctx       = stream.Context()
	)

	grpcMD, ok := grpcmd.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("key and window information are not passed in grpc metadata")
	}

	// get key from gPRC metadata
	key, err = getValueForKey(grpcMD, DatumKey)
	if err != nil {
		return err
	}

	// get window start and end time from grpc metadata
	var st, et string
	st, err = getValueForKey(grpcMD, WinStartTime)
	if err != nil {
		return err
	}

	et, err = getValueForKey(grpcMD, WinEndTime)
	if err != nil {
		return err
	}

	startTime, _ = strconv.ParseInt(st, 10, 64)
	endTime, _ = strconv.ParseInt(et, 10, 64)

	// create interval window interface using the start and end time
	iw := NewIntervalWindow(time.UnixMilli(startTime), time.UnixMilli(endTime))

	// create metadata using interval window interface
	md = NewMetadata(iw)

	var (
		datumList []*functionpb.Datum
		wg        sync.WaitGroup
	)

	wg.Add(1)
	go func() {
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
		d, err := stream.Recv()
		if err == io.EOF {
			close(reduceCh)
			break
		}
		if err != nil {
			close(reduceCh)
			// TODO: research on gRPC errors and revisit the error handler
			return err
		}
		var hd = &handlerDatum{
			value:     d.GetValue(),
			eventTime: d.GetEventTime().EventTime.AsTime(),
			watermark: d.GetWatermark().Watermark.AsTime(),
		}
		reduceCh <- hd
	}

	wg.Wait()
	return stream.SendAndClose(&functionpb.DatumList{
		Elements: datumList,
	})
}

func getValueForKey(md grpcmd.MD, key string) (string, error) {
	var value string

	keyValue := md.Get(key)

	if len(keyValue) > 1 {
		return value, fmt.Errorf("expected extactly one value for key %s in metadata but got %d values, %s", key, len(keyValue), keyValue)
	} else if len(keyValue) == 1 {
		value = keyValue[0]
	} else {
		// the length equals zero is invalid for reduce
		// since we are using a global key, and start and end time
		// cannot be empty
		return value, fmt.Errorf("expected non empty value for key %s in metadata but got an empty value", key)
	}
	return value, nil
}
