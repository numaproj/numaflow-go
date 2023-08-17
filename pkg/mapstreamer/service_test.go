package mapstreamer

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	mapstreampb "github.com/numaproj/numaflow-go/pkg/apis/proto/mapstream/v1"
)

type MapStreamFnServerTest struct {
	ctx      context.Context
	outputCh chan mapstreampb.MapStreamResponse_Result
	grpc.ServerStream
}

func NewMapStreamFnServerTest(
	ctx context.Context,
	outputCh chan mapstreampb.MapStreamResponse_Result,
) *MapStreamFnServerTest {
	return &MapStreamFnServerTest{
		ctx:      ctx,
		outputCh: outputCh,
	}
}

func (u *MapStreamFnServerTest) Send(d *mapstreampb.MapStreamResponse_Result) error {
	u.outputCh <- *d
	return nil
}

func (u *MapStreamFnServerTest) Context() context.Context {
	return u.ctx
}

type MapStreamFnServerErrTest struct {
	ctx context.Context
	grpc.ServerStream
}

func NewMapStreamFnServerErrTest(
	ctx context.Context,

) *MapStreamFnServerErrTest {
	return &MapStreamFnServerErrTest{
		ctx: ctx,
	}
}

func (u *MapStreamFnServerErrTest) Send(_ *mapstreampb.MapStreamResponse_Result) error {
	return fmt.Errorf("send error")
}

func (u *MapStreamFnServerErrTest) Context() context.Context {
	return u.ctx
}

func TestService_MapFnStream(t *testing.T) {
	tests := []struct {
		name        string
		handler     MapStreamer
		input       *mapstreampb.MapStreamRequest
		expected    []*mapstreampb.MapStreamResponse_Result
		expectedErr bool
		streamErr   bool
	}{
		{
			name: "map_stream_fn_forward_msg",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				msg := datum.Value()
				messageCh <- NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
				close(messageCh)
			}),
			input: &mapstreampb.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
			expected: []*mapstreampb.MapStreamResponse_Result{
				{
					Keys:  []string{"client_test"},
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_without_close_stream",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				msg := datum.Value()
				messageCh <- NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
			}),
			input: &mapstreampb.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
			expected: []*mapstreampb.MapStreamResponse_Result{
				{
					Keys:  []string{"client_test"},
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_forward_to_all",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				msg := datum.Value()
				messageCh <- NewMessage(msg)
				close(messageCh)
			}),
			input: &mapstreampb.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
			expected: []*mapstreampb.MapStreamResponse_Result{
				{
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_drop_msg",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				messageCh <- MessageToDrop()
				close(messageCh)
			}),
			input: &mapstreampb.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
			expected: []*mapstreampb.MapStreamResponse_Result{
				{
					Tags:  []string{DROP},
					Value: []byte{},
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_err",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				messageCh <- MessageToDrop()
				close(messageCh)
			}),
			input: &mapstreampb.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
			expected: []*mapstreampb.MapStreamResponse_Result{
				{
					Tags:  []string{DROP},
					Value: []byte{},
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				MapperStream: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
			outputCh := make(chan mapstreampb.MapStreamResponse_Result)
			result := make([]*mapstreampb.MapStreamResponse_Result, 0)

			var udfMapStreamFnStream mapstreampb.MapStream_MapStreamFnServer
			if tt.streamErr {
				udfMapStreamFnStream = NewMapStreamFnServerErrTest(ctx)
			} else {
				udfMapStreamFnStream = NewMapStreamFnServerTest(ctx, outputCh)
			}

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					result = append(result, &msg)
				}
			}()

			err := fs.MapStreamFn(tt.input, udfMapStreamFnStream)
			close(outputCh)
			wg.Wait()

			if err != nil {
				assert.True(t, tt.expectedErr, "MapStreamFn() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapStreamFn() got = %v, want %v", result, tt.expected)
			}

		})
	}
}
