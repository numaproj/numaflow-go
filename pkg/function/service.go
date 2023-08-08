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
	"google.golang.org/grpc/codes"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	functionpb "github.com/KeranYang/numaflow-go/pkg/apis/proto/function/v1"
)

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
	metadata  handlerDatumMetadata
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

func (h *handlerDatum) Metadata() DatumMetadata {
	return h.metadata
}

// handlerDatumMetadata implements the DatumMetadata interface and is used in the map and reduce handlers.
type handlerDatumMetadata struct {
	id           string
	numDelivered uint64
}

// ID returns the ID of the datum.
func (h handlerDatumMetadata) ID() string {
	return h.id
}

// NumDelivered returns the number of times the datum has been delivered.
func (h handlerDatumMetadata) NumDelivered() uint64 {
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

// Service implements the proto gen server interface and contains the map operation
// handler and the reduce operation handler.
type Service struct {
	functionpb.UnimplementedUserDefinedFunctionServer

	Mapper       MapHandler
	MapperStream MapStreamHandler
	MapperT      MapTHandler
	Reducer      ReduceHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*functionpb.ReadyResponse, error) {
	return &functionpb.ReadyResponse{Ready: true}, nil
}

// MapFn applies a function to each datum element.
func (fs *Service) MapFn(ctx context.Context, d *functionpb.DatumRequest) (*functionpb.DatumResponseList, error) {
	var hd = handlerDatum{
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
		metadata: handlerDatumMetadata{
			id:           d.GetMetadata().GetId(),
			numDelivered: d.GetMetadata().GetNumDelivered(),
		},
	}
	messages := fs.Mapper.HandleDo(ctx, d.GetKeys(), &hd)
	var elements []*functionpb.DatumResponse
	for _, m := range messages.Items() {
		elements = append(elements, &functionpb.DatumResponse{
			Keys:  m.keys,
			Value: m.value,
			Tags:  m.tags,
		})
	}
	datumList := &functionpb.DatumResponseList{
		Elements: elements,
	}
	return datumList, nil
}

// MapStreamFn applies a function to each datum element and returns a stream.
func (fs *Service) MapStreamFn(d *functionpb.DatumRequest, stream functionpb.UserDefinedFunction_MapStreamFnServer) error {
	var hd = handlerDatum{
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
		metadata: handlerDatumMetadata{
			id:           d.GetMetadata().GetId(),
			numDelivered: d.GetMetadata().GetNumDelivered(),
		},
	}
	ctx := stream.Context()
	messageCh := make(chan Message)

	done := make(chan bool)
	go func() {
		fs.MapperStream.HandleDo(ctx, d.GetKeys(), &hd, messageCh)
		done <- true
	}()
	finished := false
	for {
		select {
		case <-done:
			finished = true
		case message, ok := <-messageCh:
			if !ok {
				// Channel already closed, not closing again.
				return nil
			}
			element := &functionpb.DatumResponse{
				Keys:  message.keys,
				Value: message.value,
				Tags:  message.tags,
			}
			err := stream.Send(element)
			// the error here is returned by stream.Send() which is already a gRPC error
			if err != nil {
				// Channel may or may not be closed, as we are not sure leave it to GC.
				return err
			}
		default:
			if finished {
				close(messageCh)
				return nil
			}
		}
	}
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
// MapTFn can be used only at source vertex by source data transformer.
func (fs *Service) MapTFn(ctx context.Context, d *functionpb.DatumRequest) (*functionpb.DatumResponseList, error) {
	var hd = handlerDatum{
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
		metadata: handlerDatumMetadata{
			id:           d.GetMetadata().GetId(),
			numDelivered: d.GetMetadata().GetNumDelivered(),
		},
	}
	messageTs := fs.MapperT.HandleDo(ctx, d.GetKeys(), &hd)
	var elements []*functionpb.DatumResponse
	for _, m := range messageTs.Items() {
		elements = append(elements, &functionpb.DatumResponse{
			EventTime: &functionpb.EventTime{EventTime: timestamppb.New(m.eventTime)},
			Keys:      m.keys,
			Value:     m.value,
			Tags:      m.tags,
		})
	}
	datumList := &functionpb.DatumResponseList{
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
		statusErr := status.Errorf(codes.InvalidArgument, "keys and window information are not passed in grpc metadata")
		return statusErr
	}

	// get window start and end time from grpc metadata
	var st, et string
	st, err = getValueFromMetadata(grpcMD, WinStartTime)
	if err != nil {
		statusErr := status.Errorf(codes.InvalidArgument, err.Error())
		return statusErr
	}

	et, err = getValueFromMetadata(grpcMD, WinEndTime)
	if err != nil {
		statusErr := status.Errorf(codes.InvalidArgument, err.Error())
		return statusErr
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
		d, recvErr := stream.Recv()
		// if EOF, close all the channels
		if recvErr == io.EOF {
			closeChannels(chanMap)
			break
		}
		if recvErr != nil {
			closeChannels(chanMap)
			// the error here is returned by stream.Recv()
			// it's already a gRPC error
			return recvErr
		}
		unifiedKey := strings.Join(d.GetKeys(), Delimiter)
		var hd = &handlerDatum{
			value:     d.GetValue(),
			eventTime: d.GetEventTime().EventTime.AsTime(),
			watermark: d.GetWatermark().Watermark.AsTime(),
			metadata: handlerDatumMetadata{
				id:           d.GetMetadata().GetId(),
				numDelivered: d.GetMetadata().GetNumDelivered(),
			},
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
					sendErr := stream.Send(datumList)
					if sendErr != nil {
						// the error here is returned by stream.Send()
						// it's already a gRPC error
						return sendErr
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

func buildDatumList(messages Messages) *functionpb.DatumResponseList {
	datumList := &functionpb.DatumResponseList{}
	for _, msg := range messages {
		datumList.Elements = append(datumList.Elements, &functionpb.DatumResponse{
			Keys:  msg.keys,
			Value: msg.value,
			Tags:  msg.tags,
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
