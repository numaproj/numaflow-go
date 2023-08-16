package sink

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

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
)

type UserDefinedSink_SinkFnServerTest struct {
	ctx     context.Context
	inputCh chan *v1.SinkRequest
	rl      *v1.SinkResponseList
	grpc.ServerStream
}

func (t *UserDefinedSink_SinkFnServerTest) SendAndClose(list *v1.SinkResponseList) error {
	t.rl = list
	return nil
}

func (t *UserDefinedSink_SinkFnServerTest) Recv() (*v1.SinkRequest, error) {
	val, ok := <-t.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (t *UserDefinedSink_SinkFnServerTest) Context() context.Context {
	return t.ctx
}

func TestService_SinkFn(t *testing.T) {
	tests := []struct {
		name     string
		sh       Sinker
		input    []*v1.SinkRequest
		expected []*v1.SinkResponse
	}{
		{
			name: "sink_fn_test_success",

			input: []*v1.SinkRequest{
				{
					Id:        "one-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Id:        "two-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Id:        "three-processed",
					Keys:      []string{"sink-test"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
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
			expected: []*v1.SinkResponse{
				{
					Success: true,
					Id:      "one-processed",
					ErrMsg:  "",
				},
				{
					Success: true,
					Id:      "two-processed",
					ErrMsg:  "",
				},
				{
					Success: true,
					Id:      "three-processed",
					ErrMsg:  "",
				},
			},
		},
		{
			name: "sink_fn_test_failure",

			input: []*v1.SinkRequest{
				{
					Id:        "one-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Id:        "two-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Id:        "three-processed",
					Keys:      []string{"sink-test-1", "sink-test-2"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
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
			expected: []*v1.SinkResponse{
				{
					Success: false,
					Id:      "one-processed",
					ErrMsg:  "unknown error",
				},
				{
					Success: false,
					Id:      "two-processed",
					ErrMsg:  "unknown error",
				},
				{
					Success: false,
					Id:      "three-processed",
					ErrMsg:  "unknown error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := Service{
				Sinker: tt.sh,
			}
			ich := make(chan *v1.SinkRequest)
			udfReduceFnStream := &UserDefinedSink_SinkFnServerTest{
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

			if !reflect.DeepEqual(udfReduceFnStream.rl.Responses, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", udfReduceFnStream.rl.Responses, tt.expected)
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
		want        *v1.ReadyResponse
		wantErr     bool
	}{
		{
			name: "is_ready",
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &v1.ReadyResponse{
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
