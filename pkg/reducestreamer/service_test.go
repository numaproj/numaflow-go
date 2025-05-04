package reducestreamer

import (
	"context"
	"io"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	reducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

type ReduceFnServerTest struct {
	ctx      context.Context
	inputCh  chan *reducepb.ReduceRequest
	outputCh chan *reducepb.ReduceResponse
	grpc.ServerStream
}

func NewReduceFnServerTest(ctx context.Context,
	inputCh chan *reducepb.ReduceRequest,
	outputCh chan *reducepb.ReduceResponse) *ReduceFnServerTest {
	return &ReduceFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *ReduceFnServerTest) Send(list *reducepb.ReduceResponse) error {
	u.outputCh <- list
	return nil
}

func (u *ReduceFnServerTest) Recv() (*reducepb.ReduceRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *ReduceFnServerTest) Context() context.Context {
	return u.ctx
}

func TestService_ReduceFn(t *testing.T) {

	tests := []struct {
		name        string
		handler     func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata)
		input       []*reducepb.ReduceRequest
		expected    []*reducepb.ReduceResponse
		expectedErr bool
	}{
		{
			name: "reduce_fn_forward_msg_same_keys",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"})
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client_test"},
						Value: []byte(strconv.Itoa(60)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_multiple_keys",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"})
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client3"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client3"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client3_test"},
						Value: []byte(strconv.Itoa(60)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_forward_to_all",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum)))
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Value: []byte(strconv.Itoa(60)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
			},
		},
		{
			name: "reduce_fn_forward_msg_drop_msg",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- MessageToDrop()
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Tags:  []string{DROP},
						Value: []byte{},
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
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
				creatorHandle: SimpleCreatorWithReduceStreamFn(tt.handler),
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{winStartTime: "60000", winEndTime: "120000"}))

			inputCh := make(chan *reducepb.ReduceRequest)
			outputCh := make(chan *reducepb.ReduceResponse)
			result := make([]*reducepb.ReduceResponse, 0)

			udfReduceFnStream := NewReduceFnServerTest(ctx, inputCh, outputCh)

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.ReduceFn(udfReduceFnStream)
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
				t.Errorf("ReduceFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestService_ReduceFn_CloseTask(t *testing.T) {
	tests := []struct {
		name        string
		handler     func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata)
		input       []*reducepb.ReduceRequest
		expected    []*reducepb.ReduceResponse
		expectedErr bool
	}{
		{
			name: "reduce_fn_with_explicit_close",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"})
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(40)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				// Explicit CLOSE operation for the window
				{
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_CLOSE,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(30)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(70)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: true,
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_with_slow_processing",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					// Simulate slow processing
					time.Sleep(50 * time.Millisecond)
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"})
			},
			input: []*reducepb.ReduceRequest{
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_APPEND,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				// Explicit CLOSE operation for the window
				{
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_CLOSE,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(30)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: true,
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_with_multiple_windows",
			handler: func(ctx context.Context, keys []string, rch <-chan Datum, och chan<- Message, md Metadata) {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				och <- NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"})
			},
			input: []*reducepb.ReduceRequest{
				// Window 1
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client1"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				// Window 2
				{
					Payload: &reducepb.ReduceRequest_Payload{
						Keys:      []string{"client2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_OPEN,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(120000)),
								End:   timestamppb.New(time.UnixMilli(180000)),
								Slot:  "slot-1",
							},
						},
					},
				},
				// Close Window 1
				{
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_CLOSE,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(60000)),
								End:   timestamppb.New(time.UnixMilli(120000)),
								Slot:  "slot-0",
							},
						},
					},
				},
				// Close Window 2
				{
					Operation: &reducepb.ReduceRequest_WindowOperation{
						Event: reducepb.ReduceRequest_WindowOperation_CLOSE,
						Windows: []*reducepb.Window{
							{
								Start: timestamppb.New(time.UnixMilli(120000)),
								End:   timestamppb.New(time.UnixMilli(180000)),
								Slot:  "slot-1",
							},
						},
					},
				},
			},
			expected: []*reducepb.ReduceResponse{
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(10)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: false,
				},
				{
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(60000)),
						End:   timestamppb.New(time.UnixMilli(120000)),
						Slot:  "slot-0",
					},
					EOF: true,
				},
				{
					Result: &reducepb.ReduceResponse_Result{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(120000)),
						End:   timestamppb.New(time.UnixMilli(180000)),
						Slot:  "slot-1",
					},
					EOF: false,
				},
				{
					Window: &reducepb.Window{
						Start: timestamppb.New(time.UnixMilli(120000)),
						End:   timestamppb.New(time.UnixMilli(180000)),
						Slot:  "slot-1",
					},
					EOF: true,
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				creatorHandle: SimpleCreatorWithReduceStreamFn(tt.handler),
				shutdownCh:    make(chan struct{}, 1),
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{winStartTime: "60000", winEndTime: "120000"}))

			inputCh := make(chan *reducepb.ReduceRequest)
			outputCh := make(chan *reducepb.ReduceResponse)
			result := make([]*reducepb.ReduceResponse, 0)

			udfReduceFnStream := NewReduceFnServerTest(ctx, inputCh, outputCh)

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.ReduceFn(udfReduceFnStream)
				close(outputCh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					result = append(result, msg)
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

			// Sort results for comparison
			sort.Slice(result, func(i, j int) bool {
				// First sort by EOF (non-EOF messages first)
				if result[i].EOF != result[j].EOF {
					return !result[i].EOF
				}

				// Then sort by window start time
				iStart := result[i].Window.Start.AsTime().UnixMilli()
				jStart := result[j].Window.Start.AsTime().UnixMilli()
				if iStart != jStart {
					return iStart < jStart
				}

				// Then sort by value if present
				if result[i].Result != nil && result[j].Result != nil {
					return string(result[i].Result.Value) < string(result[j].Result.Value)
				}

				// EOF messages without results come last
				return result[i].Result != nil
			})

			// Sort expected results the same way
			sort.Slice(tt.expected, func(i, j int) bool {
				// First sort by EOF (non-EOF messages first)
				if tt.expected[i].EOF != tt.expected[j].EOF {
					return !tt.expected[i].EOF
				}

				// Then sort by window start time
				iStart := tt.expected[i].Window.Start.AsTime().UnixMilli()
				jStart := tt.expected[j].Window.Start.AsTime().UnixMilli()
				if iStart != jStart {
					return iStart < jStart
				}

				// Then sort by value if present
				if tt.expected[i].Result != nil && tt.expected[j].Result != nil {
					return string(tt.expected[i].Result.Value) < string(tt.expected[j].Result.Value)
				}

				// EOF messages without results come last
				return tt.expected[i].Result != nil
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
