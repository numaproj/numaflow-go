package smap

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/smapfn"
	"github.com/numaproj/numaflow-go/pkg/function"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MapStreamFnServerTest struct {
	ctx      context.Context
	outputCh chan smapfn.MapStreamResponse
	grpc.ServerStream
}

func NewMapStreamFnServerTest(
	ctx context.Context,
	outputCh chan smapfn.MapStreamResponse,
) *MapStreamFnServerTest {
	return &MapStreamFnServerTest{
		ctx:      ctx,
		outputCh: outputCh,
	}
}

func (u *MapStreamFnServerTest) Send(d *smapfn.MapStreamResponse) error {
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

func (u *MapStreamFnServerErrTest) Send(_ *smapfn.MapStreamResponse) error {
	return fmt.Errorf("send error")
}

func (u *MapStreamFnServerErrTest) Context() context.Context {
	return u.ctx
}

func TestService_MapFnStream(t *testing.T) {
	tests := []struct {
		name        string
		handler     function.MapStreamHandler
		input       *smapfn.MapStreamRequest
		expected    []*smapfn.MapStreamResponse
		expectedErr bool
		streamErr   bool
	}{
		{
			name: "map_stream_fn_forward_msg",
			handler: function.MapStreamFunc(func(ctx context.Context, keys []string, datum function.Datum, messageCh chan<- function.Message) {
				msg := datum.Value()
				messageCh <- function.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
				close(messageCh)
			}),
			input: &smapfn.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: &smapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
				Watermark: &smapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
			},
			expected: []*smapfn.MapStreamResponse{
				{
					Keys:  []string{"client_test"},
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_without_close_stream",
			handler: function.MapStreamFunc(func(ctx context.Context, keys []string, datum function.Datum, messageCh chan<- function.Message) {
				msg := datum.Value()
				messageCh <- function.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
			}),
			input: &smapfn.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: &smapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
				Watermark: &smapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
			},
			expected: []*smapfn.MapStreamResponse{
				{
					Keys:  []string{"client_test"},
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_forward_to_all",
			handler: function.MapStreamFunc(func(ctx context.Context, keys []string, datum function.Datum, messageCh chan<- function.Message) {
				msg := datum.Value()
				messageCh <- function.NewMessage(msg)
				close(messageCh)
			}),
			input: &smapfn.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: &smapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
				Watermark: &smapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
			},
			expected: []*smapfn.MapStreamResponse{
				{
					Value: []byte(`test`),
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_msg_drop_msg",
			handler: function.MapStreamFunc(func(ctx context.Context, keys []string, datum function.Datum, messageCh chan<- function.Message) {
				messageCh <- function.MessageToDrop()
				close(messageCh)
			}),
			input: &smapfn.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: &smapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
				Watermark: &smapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
			},
			expected: []*smapfn.MapStreamResponse{
				{
					Tags:  []string{function.DROP},
					Value: []byte{},
				},
			},
			expectedErr: false,
		},
		{
			name: "map_stream_fn_forward_err",
			handler: function.MapStreamFunc(func(ctx context.Context, keys []string, datum function.Datum, messageCh chan<- function.Message) {
				messageCh <- function.MessageToDrop()
				close(messageCh)
			}),
			input: &smapfn.MapStreamRequest{
				Keys:      []string{"client"},
				Value:     []byte(`test`),
				EventTime: &smapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
				Watermark: &smapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
			},
			expected: []*smapfn.MapStreamResponse{
				{
					Tags:  []string{function.DROP},
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
			outputCh := make(chan smapfn.MapStreamResponse)
			result := make([]*smapfn.MapStreamResponse, 0)

			var udfMapStreamFnStream smapfn.MapStream_MapStreamFnServer
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
