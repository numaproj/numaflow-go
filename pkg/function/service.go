package function

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
)

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct {
	id           string
	value        []byte
	eventTime    time.Time
	watermark    time.Time
	numDelivered uint64
}

func (h *handlerDatum) ID() string {
	return h.id
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

func (h *handlerDatum) NumDelivered() uint64 {
	return h.numDelivered
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

// MapFn applies a function to each datum element.
func (fs *Service) MapFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	var hd = handlerDatum{
		id:           d.Id,
		value:        d.GetValue(),
		eventTime:    d.GetEventTime().EventTime.AsTime(),
		watermark:    d.GetWatermark().Watermark.AsTime(),
		numDelivered: d.NumDelivered,
	}
	messages := fs.Mapper.HandleDo(ctx, d.GetKeys(), &hd)
	var elements []*functionpb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.Datum{
			Keys:  m.keys,
			Value: m.value,
		})
	}
	datumList := &functionpb.DatumList{
		Elements: elements,
	}
	return datumList, nil
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
// MapTFn can be used only at source vertex by source data transformer.
func (fs *Service) MapTFn(ctx context.Context, d *functionpb.Datum) (*functionpb.DatumList, error) {
	var hd = handlerDatum{
		id:           d.Id,
		value:        d.GetValue(),
		eventTime:    d.GetEventTime().EventTime.AsTime(),
		watermark:    d.GetWatermark().Watermark.AsTime(),
		numDelivered: d.NumDelivered,
	}
	messageTs := fs.MapperT.HandleDo(ctx, d.GetKeys(), &hd)
	var elements []*functionpb.Datum
	for _, m := range messageTs.Items() {
		elements = append(elements, &functionpb.Datum{
			EventTime: &functionpb.EventTime{EventTime: timestamppb.New(m.eventTime)},
			Keys:      m.keys,
			Value:     m.value,
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
		md        Metadata
		err       error
		startTime int64
		endTime   int64
		ctx       = stream.Context()
		chanMap   = make(map[string]chan Datum)
		mu        sync.RWMutex
		g         errgroup.Group
	)

	grpcMD, ok := grpcmd.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("keys and window information are not passed in grpc metadata")
	}

	// get window start and end time from grpc metadata
	var st, et string
	st, err = getValueFromMetadata(grpcMD, WinStartTime)
	if err != nil {
		return err
	}

	et, err = getValueFromMetadata(grpcMD, WinEndTime)
	if err != nil {
		return err
	}

	startTime, _ = strconv.ParseInt(st, 10, 64)
	endTime, _ = strconv.ParseInt(et, 10, 64)

	// create interval window interface using the start and end time
	iw := NewIntervalWindow(time.UnixMilli(startTime), time.UnixMilli(endTime))

	// create metadata using interval window interface
	md = NewMetadata(iw)

	// read messages from the stream and write the messages to corresponding channels
	// if the channel is not created, create the channel and invoke the HandleDo
	for {
		d, err := stream.Recv()
		// if EOF, close all the channels
		if err == io.EOF {
			closeChannels(chanMap)
			break
		}
		if err != nil {
			closeChannels(chanMap)
			// TODO: research on gRPC errors and revisit the error handler
			return err
		}
		unifiedKey := strings.Join(d.GetKeys(), Delimiter)
		var hd = &handlerDatum{
			id:           d.Id,
			value:        d.GetValue(),
			eventTime:    d.GetEventTime().EventTime.AsTime(),
			watermark:    d.GetWatermark().Watermark.AsTime(),
			numDelivered: d.GetNumDelivered(),
		}

		ch, chok := chanMap[unifiedKey]
		if !chok {
			ch = make(chan Datum)
			chanMap[unifiedKey] = ch

			func(k []string, ch chan Datum) {
				g.Go(func() error {
					// we stream the messages to the user by writing messages
					// to channel and wait until we get the result and stream
					// the result back to the client (numaflow).
					messages := fs.Reducer.HandleDo(ctx, k, ch, md)
					datumList := buildDatumList(messages)

					// stream.Send() is not thread safe.
					mu.Lock()
					defer mu.Unlock()
					err := stream.Send(datumList)
					if err != nil {
						// TODO: research on gRPC errors and revisit the error handler
						return err
					}
					return nil
				})
			}(d.GetKeys(), ch)
		}
		ch <- hd
	}

	// wait until all the HandleDo return
	return g.Wait()
}

func buildDatumList(messages Messages) *functionpb.DatumList {
	datumList := &functionpb.DatumList{}
	for _, msg := range messages {
		datumList.Elements = append(datumList.Elements, &functionpb.Datum{
			Keys:  msg.keys,
			Value: msg.value,
		})
	}

	return datumList
}

func closeChannels(chanMap map[string]chan Datum) {
	for _, ch := range chanMap {
		close(ch)
	}
}

func getValueFromMetadata(md grpcmd.MD, k string) (string, error) {
	var value string

	keyValue := md.Get(k)

	if len(keyValue) > 1 {
		return value, fmt.Errorf("expected extactly one value for keys %s in metadata but got %d values, %s", k, len(keyValue), keyValue)
	} else if len(keyValue) == 1 {
		value = keyValue[0]
	} else {
		// the length equals zero is invalid for reduce
		// since we are using a global keys, and start and end time
		// cannot be empty
		return value, fmt.Errorf("expected non empty value for keys %s in metadata but got an empty value", k)
	}
	return value, nil
}
