package sessionreducer

import (
	"context"
	"io"
	"log"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	sessionreducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1"
)

type SessionReduceFnServerTest struct {
	ctx      context.Context
	inputCh  chan *sessionreducepb.SessionReduceRequest
	outputCh chan *sessionreducepb.SessionReduceResponse
	grpc.ServerStream
}

func NewReduceFnServerTest(ctx context.Context,
	inputCh chan *sessionreducepb.SessionReduceRequest,
	outputCh chan *sessionreducepb.SessionReduceResponse) *SessionReduceFnServerTest {
	return &SessionReduceFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *SessionReduceFnServerTest) Send(list *sessionreducepb.SessionReduceResponse) error {
	u.outputCh <- list
	return nil
}

func (u *SessionReduceFnServerTest) Recv() (*sessionreducepb.SessionReduceRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *SessionReduceFnServerTest) Context() context.Context {
	return u.ctx
}

type SessionSum struct {
	sum *atomic.Int32
}

func (s *SessionSum) SessionReduce(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message) {
	for val := range inputCh {
		msgVal, _ := strconv.Atoi(string(val.Value()))
		s.sum.Add(int32(msgVal))
	}
	outputCh <- NewMessage([]byte(strconv.Itoa(int(s.sum.Load())))).WithKeys([]string{keys[0] + "_test"})
}

func (s *SessionSum) Accumulator(ctx context.Context) []byte {
	return []byte(strconv.Itoa(int(s.sum.Load())))
}

func (s *SessionSum) MergeAccumulator(ctx context.Context, accumulator []byte) {
	val, err := strconv.Atoi(string(accumulator))
	if err != nil {
		log.Println("unable to convert the accumulator value to int: ", err.Error())
		return
	}
	s.sum.Add(int32(val))
}

type SessionSumCreator struct {
}

func (s *SessionSumCreator) Create() SessionReducer {
	return NewSessionSum()
}

func NewSessionSum() SessionReducer {
	return &SessionSum{
		sum: atomic.NewInt32(0),
	}
}

func TestService_SessionReduceFn(t *testing.T) {

	tests := []struct {
		name        string
		handler     SessionReducerCreator
		input       []*sessionreducepb.SessionReduceRequest
		expected    []*sessionreducepb.SessionReduceResponse
		expectedErr bool
	}{
		{
			name:    "open_append_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-1"},
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
								Keys:  []string{"client"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
								Keys:  []string{"client"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
								Keys:  []string{"client"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
								Keys:  []string{"client"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
								Keys:  []string{"client"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client_test"},
						Value: []byte(strconv.Itoa(60)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
						Keys:  []string{"client"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_expand_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(75000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(79000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(75000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(79000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(75000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(79000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_merge_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(85000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(88000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_expand_append_merge_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(95000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(80)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(98000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_merge_append_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_APPEND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(30)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(85000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(50)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(88000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_merge_expand_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_EXPAND,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(95000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(98000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(30)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(95000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(50)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(98000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name:    "open_merge_merge_close",
			handler: &SessionSumCreator{},
			input: []*sessionreducepb.SessionReduceRequest{
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(75000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(70000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(78000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(80000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Payload: &sessionreducepb.SessionReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_OPEN,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(80000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(80000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_MERGE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(80000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
				{
					Operation: &sessionreducepb.SessionReduceRequest_WindowOperation{
						Event: sessionreducepb.SessionReduceRequest_WindowOperation_CLOSE,
						KeyedWindows: []*sessionreducepb.KeyedWindow{
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(85000)),
								Slot:  "slot-0",
								Keys:  []string{"client1"},
							},
							{
								Start: timestamppb.New(time.UnixMilli(50000)),
								End:   timestamppb.New(time.UnixMilli(88000)),
								Slot:  "slot-0",
								Keys:  []string{"client2"},
							},
						},
					},
				},
			},
			expected: []*sessionreducepb.SessionReduceResponse{
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(30)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(50000)),
						End:   timestamppb.New(time.UnixMilli(85000)),
						Slot:  "slot-0",
						Keys:  []string{"client1"},
					},
					EOF: false,
				},
				{
					Result: &sessionreducepb.SessionReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(50)),
					},
					KeyedWindow: &sessionreducepb.KeyedWindow{
						Start: timestamppb.New(time.UnixMilli(50000)),
						End:   timestamppb.New(time.UnixMilli(88000)),
						Slot:  "slot-0",
						Keys:  []string{"client2"},
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				creatorHandle: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.

			inputCh := make(chan *sessionreducepb.SessionReduceRequest)
			outputCh := make(chan *sessionreducepb.SessionReduceResponse)
			result := make([]*sessionreducepb.SessionReduceResponse, 0)

			udfReduceFnStream := NewReduceFnServerTest(context.Background(), inputCh, outputCh)

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.SessionReduceFn(udfReduceFnStream)
				close(outputCh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					if !msg.EOF {
						result = append(result, msg)
					}
				}
			}()

			for _, val := range tt.input {
				udfReduceFnStream.inputCh <- val
			}
			close(udfReduceFnStream.inputCh)
			wg.Wait()

			if (err != nil) != tt.expectedErr {
				t.Errorf("ReduceFn() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}

			//sort and compare, since order of the output doesn't matter
			sort.Slice(result, func(i, j int) bool {
				return string(result[i].Result.Value) < string(result[j].Result.Value)
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SessionReduceFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
