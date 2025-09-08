package sinker

import (
	"context"
	"io"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
)

type SinkFnServerTest struct {
	ctx     context.Context
	inputCh chan *sinkpb.SinkRequest
	rl      []*sinkpb.SinkResponse
	grpc.ServerStream
}

func (t *SinkFnServerTest) Send(response *sinkpb.SinkResponse) error {
	t.rl = append(t.rl, response)
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
		expected []*sinkpb.SinkResponse
	}{
		{
			name: "sink_fn_test_success",

			input: []*sinkpb.SinkRequest{
				{
					Handshake: &sinkpb.Handshake{
						Sot: true,
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "one-processed",
						Keys:      []string{"sink-test"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-1"},
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "two-processed",
						Keys:      []string{"sink-test"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-2"},
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "three-processed",
						Keys:      []string{"sink-test"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-3"},
					},
				},
				{
					Status: &sinkpb.TransmissionStatus{Eot: true},
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
			expected: []*sinkpb.SinkResponse{
				{
					Handshake: &sinkpb.Handshake{
						Sot: true,
					},
				},
				{
					Results: []*sinkpb.SinkResponse_Result{
						{
							Status:   sinkpb.Status_SUCCESS,
							Id:       "one-processed",
							ErrMsg:   "",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
						{
							Status:   sinkpb.Status_SUCCESS,
							Id:       "two-processed",
							ErrMsg:   "",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
						{
							Status:   sinkpb.Status_SUCCESS,
							Id:       "three-processed",
							ErrMsg:   "",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
					},
				},
				{
					Status: &sinkpb.TransmissionStatus{Eot: true},
				},
			},
		},
		{
			name: "sink_fn_test_failure",
			input: []*sinkpb.SinkRequest{
				{
					Handshake: &sinkpb.Handshake{
						Sot: true,
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "one-processed",
						Keys:      []string{"sink-test-1", "sink-test-2"},
						Value:     []byte(strconv.Itoa(10)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-1"},
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "two-processed",
						Keys:      []string{"sink-test-1", "sink-test-2"},
						Value:     []byte(strconv.Itoa(20)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
				},
				{
					Request: &sinkpb.SinkRequest_Request{
						Id:        "three-processed",
						Keys:      []string{"sink-test-1", "sink-test-2"},
						Value:     []byte(strconv.Itoa(30)),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
						Headers:   map[string]string{"x-txn-id": "test-txn-2"},
					},
				},
				{
					Status: &sinkpb.TransmissionStatus{Eot: true},
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
			expected: []*sinkpb.SinkResponse{
				{
					Handshake: &sinkpb.Handshake{
						Sot: true,
					},
				},
				{
					Results: []*sinkpb.SinkResponse_Result{
						{
							Status:   sinkpb.Status_FAILURE,
							Id:       "one-processed",
							ErrMsg:   "unknown error",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
						{
							Status:   sinkpb.Status_FAILURE,
							Id:       "two-processed",
							ErrMsg:   "unknown error",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
						{
							Status:   sinkpb.Status_FAILURE,
							Id:       "three-processed",
							ErrMsg:   "unknown error",
							Metadata: &common.Metadata{SysMetadata: map[string]*common.KeyValueGroup{}, UserMetadata: map[string]*common.KeyValueGroup{}},
						},
					},
				},
				{
					Status: &sinkpb.TransmissionStatus{Eot: true},
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
				err := ss.SinkFn(udfReduceFnStream)
				assert.NoError(t, err)
			}()

			for _, val := range tt.input {
				ich <- val
			}
			close(ich)

			wg.Wait()

			for i, val := range tt.expected {
				assert.Equal(t, val, udfReduceFnStream.rl[i])
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
