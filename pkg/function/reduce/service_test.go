package reduce

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

	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/reducefn"
	"github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/util"
)

type ReduceFnServerTest struct {
	ctx      context.Context
	inputCh  chan *reducefn.ReduceRequest
	outputCh chan *reducefn.ReduceResponseList
	grpc.ServerStream
}

func NewReduceFnServerTest(ctx context.Context,
	inputCh chan *reducefn.ReduceRequest,
	outputCh chan *reducefn.ReduceResponseList) *ReduceFnServerTest {
	return &ReduceFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *ReduceFnServerTest) Send(list *reducefn.ReduceResponseList) error {
	u.outputCh <- list
	return nil
}

func (u *ReduceFnServerTest) Recv() (*reducefn.ReduceRequest, error) {
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
		handler     function.ReduceHandler
		input       []*reducefn.ReduceRequest
		expected    *reducefn.ReduceResponseList
		expectedErr bool
	}{
		{
			name: "reduce_fn_forward_msg_same_keys",
			handler: function.ReduceFunc(func(ctx context.Context, keys []string, rch <-chan function.Datum, md function.Metadata) function.Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return function.MessagesBuilder().Append(function.NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
			}),
			input: []*reducefn.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &reducefn.ReduceResponseList{
				Elements: []*reducefn.ReduceResponse{
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
			handler: function.ReduceFunc(func(ctx context.Context, keys []string, rch <-chan function.Datum, md function.Metadata) function.Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return function.MessagesBuilder().Append(function.NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
			}),
			input: []*reducefn.ReduceRequest{
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &reducefn.ReduceResponseList{
				Elements: []*reducefn.ReduceResponse{
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
			handler: function.ReduceFunc(func(ctx context.Context, keys []string, rch <-chan function.Datum, md function.Metadata) function.Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return function.MessagesBuilder().Append(function.NewMessage([]byte(strconv.Itoa(sum))))
			}),
			input: []*reducefn.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &reducefn.ReduceResponseList{
				Elements: []*reducefn.ReduceResponse{
					{
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_drop_msg",
			handler: function.ReduceFunc(func(ctx context.Context, keys []string, rch <-chan function.Datum, md function.Metadata) function.Messages {
				sum := 0
				for val := range rch {
					msgVal, _ := strconv.Atoi(string(val.Value()))
					sum += msgVal
				}
				return function.MessagesBuilder().Append(function.MessageToDrop())
			}),
			input: []*reducefn.ReduceRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &common.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &common.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &reducefn.ReduceResponseList{
				Elements: []*reducefn.ReduceResponse{
					{
						Tags:  []string{function.DROP},
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
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{util.WinStartTime: "60000", util.WinEndTime: "120000"}))

			inputCh := make(chan *reducefn.ReduceRequest)
			outputCh := make(chan *reducefn.ReduceResponseList)
			result := &reducefn.ReduceResponseList{}

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
					result.Elements = append(result.Elements, msg.Elements...)
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
			sort.Slice(result.Elements, func(i, j int) bool {
				return string(result.Elements[i].Value) < string(result.Elements[j].Value)
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
