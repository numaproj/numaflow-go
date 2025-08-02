package reducestreamer

import (
	"context"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
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

func TestService_ReduceFn_PanicHandling(t *testing.T) {
	// Test handler that panics during reduce stream operation
	panicHandler := func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata) {
		// Read the first message, then panic
		<-inputCh
		panic("test panic in reduce streamer handler")
	}

	shutdownCh := make(chan struct{}, 1)

	fs := &Service{
		creatorHandle: SimpleCreatorWithReduceStreamFn(panicHandler),
		shutdownCh:    shutdownCh,
	}

	ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{winStartTime: "60000", winEndTime: "120000"}))

	inputCh := make(chan *reducepb.ReduceRequest)
	outputCh := make(chan *reducepb.ReduceResponse)

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
		// Consume any messages from output channel
		for range outputCh {
			// Just drain the channel
		}
	}()

	// Send a request that will trigger the panic
	inputCh <- &reducepb.ReduceRequest{
		Payload: &reducepb.ReduceRequest_Payload{
			Keys:      []string{"test-key"},
			Value:     []byte("test-value"),
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
	}

	// Give a small delay to ensure error propagation completes
	time.Sleep(100 * time.Millisecond)
	// Close input to trigger EOF
	close(inputCh)
	wg.Wait()

	// Verify that an error was returned
	if err == nil {
		t.Error("Expected error due to panic, but got nil")
		return
	}

	// Verify that the error contains panic information
	if !strings.Contains(err.Error(), "UDF_EXECUTION_ERROR") {
		t.Errorf("Expected error to contain 'UDF_EXECUTION_ERROR', got: %v", err)
	}

	// Verify that shutdown channel was triggered
	select {
	case <-shutdownCh:
		// Expected - shutdown was triggered
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected shutdown channel to be triggered, but it wasn't")
	}
}
