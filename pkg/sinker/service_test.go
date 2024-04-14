package sinker

import (
	"context"
	"io"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
)

type SinkFnServerTest struct {
	ctx     context.Context
	inputCh chan *sinkpb.SinkRequest
	rl      *sinkpb.SinkResponse
	grpc.ServerStream
}

func (t *SinkFnServerTest) SendAndClose(list *sinkpb.SinkResponse) error {
	t.rl = list
	return nil
}

func (t *SinkFnServerTest) Recv() (*sinkpb.SinkRequest, error) {
	val, ok := <-t.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (t *SinkFnServerTest) Context() context.Context {
	return t.ctx
}

func TestService_SinkFn(t *testing.T) {
	tests := []struct {
		name     string
		sh       Sinker
		input    []*sinkpb.SinkRequest
		expected []*sinkpb.SinkResponse_Result
	}{
		{
			name: "sink_fn_test_success",

			input: []*sinkpb.SinkRequest{
				{
					Id:        "one-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
					Headers:   map[string]string{"x-txn-id": "test-txn-1"},
				},
				{
					Id:        "two-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
					Headers:   map[string]string{"x-txn-id": "test-txn-2"},
				},
				{
					Id:        "three-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
					Headers:   map[string]string{"x-txn-id": "test-txn-3"},
				},
			},
			sh: SinkerFunc(func(ctx context.Context, rch <-chan Datum) Responses {
				result := ResponsesBuilder()
				for d := range rch {
					id := d.ID()
					result = result.Append(ResponseOK(id))
				}
				return result
			}),
			expected: []*sinkpb.SinkResponse_Result{
				{
					Status: sinkpb.Status_SUCCESS,
					Id:     "one-processed",
					ErrMsg: "",
				},
				{
					Status: sinkpb.Status_SUCCESS,
					Id:     "two-processed",
					ErrMsg: "",
				},
				{
					Status: sinkpb.Status_SUCCESS,
					Id:     "three-processed",
					ErrMsg: "",
				},
			},
		},
		{
			name: "sink_fn_test_failure",

			input: []*sinkpb.SinkRequest{
				{
					Id:        "one-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
					Headers:   map[string]string{"x-txn-id": "test-txn-1"},
				},
				{
					Id:        "two-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
				{
					Id:        "three-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
					Headers:   map[string]string{"x-txn-id": "test-txn-2"},
				},
			},
			sh: SinkerFunc(func(ctx context.Context, rch <-chan Datum) Responses {
				result := ResponsesBuilder()
				for d := range rch {
					id := d.ID()
					result = result.Append(ResponseFailure(id, "unknown error"))
				}
				return result
			}),
			expected: []*sinkpb.SinkResponse_Result{
				{
					Status: sinkpb.Status_FAILURE,
					Id:     "one-processed",
					ErrMsg: "unknown error",
				},
				{
					Status: sinkpb.Status_FAILURE,
					Id:     "two-processed",
					ErrMsg: "unknown error",
				},
				{
					Status: sinkpb.Status_FAILURE,
					Id:     "three-processed",
					ErrMsg: "unknown error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := Service{
				Sinker: tt.sh,
			}
			ich := make(chan *sinkpb.SinkRequest)
			udfReduceFnStream := &SinkFnServerTest{
				ctx:     context.Background(),
				inputCh: ich,
			}

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = ss.SinkFn(udfReduceFnStream)
			}()

			for _, val := range tt.input {
				ich <- val
			}
			close(ich)

			wg.Wait()

			if !reflect.DeepEqual(udfReduceFnStream.rl.Results, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", udfReduceFnStream.rl.Results, tt.expected)
			}
		})
	}
}

func TestService_IsReady(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name        string
		sinkHandler Sinker
		args        args
		want        *sinkpb.ReadyResponse
		wantErr     bool
	}{
		{
			name: "is_ready",
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &sinkpb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Sinker: tt.sinkHandler,
			}
			got, err := fs.IsReady(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReady() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsReady() got = %v, want %v", got, tt.want)
			}
		})
	}
}
