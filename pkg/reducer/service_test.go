package reducer

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
	"github.com/numaproj/numaflow-go/pkg/shared"
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
		handler     Reducer
		input       []*reducepb.ReduceRequest
		expected    *reducepb.ReduceResponse
		expectedErr bool
	}{
		{
			name: "reduce_fn_forward_msg_same_keys",
			handler: ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
			}),
			input: []*reducepb.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			expected: &reducepb.ReduceResponse{
				Results: []*reducepb.ReduceResponse_Result{
					{
						Keys:  []string{"client_test"},
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_multiple_keys",
			handler: ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
			}),
			input: []*reducepb.ReduceRequest{
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			expected: &reducepb.ReduceResponse{
				Results: []*reducepb.ReduceResponse_Result{
					{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					{
						Keys:  []string{"client3_test"},
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_forward_to_all",
			handler: ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))))
			}),
			input: []*reducepb.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			expected: &reducepb.ReduceResponse{
				Results: []*reducepb.ReduceResponse_Result{
					{
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_drop_msg",
			handler: ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return MessagesBuilder().Append(MessageToDrop())
			}),
			input: []*reducepb.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			expected: &reducepb.ReduceResponse{
				Results: []*reducepb.ReduceResponse_Result{
					{
						Tags:  []string{DROP},
						Value: []byte{},
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Reducer: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{shared.WinStartTime: "60000", shared.WinEndTime: "120000"}))

			inputCh := make(chan *reducepb.ReduceRequest)
			outputCh := make(chan *reducepb.ReduceResponse)
			result := &reducepb.ReduceResponse{}

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
					result.Results = append(result.Results, msg.Results...)
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
			sort.Slice(result.Results, func(i, j int) bool {
				return string(result.Results[i].Value) < string(result.Results[j].Value)
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
