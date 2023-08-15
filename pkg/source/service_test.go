package source

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/source/model"
)

type fields struct {
	pendingHandler PendingHandler
	ackHandler     AckHandler
	readHandler    ReadHandler
}

func TestService_IsReady(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sourcepb.ReadyResponse
		wantErr bool
	}{
		{
			name: "is_ready",
			fields: fields{
				pendingHandler: nil,
				ackHandler:     nil,
				readHandler:    nil,
			},
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &sourcepb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				PendingHandler: tt.fields.pendingHandler,
				AckHandler:     tt.fields.ackHandler,
				ReadHandler:    tt.fields.readHandler,
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

type ReadFnServerTest struct {
	ctx      context.Context
	outputCh chan sourcepb.DatumResponse
	grpc.ServerStream
}

func NewReadFnServerTest(
	ctx context.Context,
	outputCh chan sourcepb.DatumResponse,
) *ReadFnServerTest {
	return &ReadFnServerTest{
		ctx:      ctx,
		outputCh: outputCh,
	}
}

func (t *ReadFnServerTest) Send(d *sourcepb.DatumResponse) error {
	t.outputCh <- *d
	return nil
}

func (t *ReadFnServerTest) Context() context.Context {
	return t.ctx
}

type ReadFnServerErrTest struct {
	ctx context.Context
	grpc.ServerStream
}

func NewReadFnServerErrTest(
	ctx context.Context,
) *ReadFnServerErrTest {
	return &ReadFnServerErrTest{
		ctx: ctx,
	}
}

func (te *ReadFnServerErrTest) Send(_ *sourcepb.DatumResponse) error {
	return fmt.Errorf("send error")
}

func (te *ReadFnServerErrTest) Context() context.Context {
	return te.ctx
}

func TestService_ReadFn(t *testing.T) {
	testEventTime := time.Now()
	tests := []struct {
		name        string
		fields      fields
		input       *sourcepb.ReadRequest
		expected    []*sourcepb.DatumResponse
		expectedErr bool
	}{
		{
			name: "read_fn_read_msg",
			fields: fields{
				readHandler: ReadFunc(func(ctx context.Context, request ReadRequest, messageCh chan<- model.Message) {
					msg := model.NewMessage([]byte(`test`), model.Offset{}, testEventTime)
					messageCh <- msg.WithKeys([]string{"test-key"})
					close(messageCh)
				}),
			},
			input: &sourcepb.ReadRequest{
				Request: &sourcepb.ReadRequest_Request{
					NumRecords:  1,
					TimeoutInMs: 1000,
				},
			},
			expected: []*sourcepb.DatumResponse{
				{
					Result: &sourcepb.DatumResponse_Result{
						Payload:   []byte(`test`),
						Offset:    &sourcepb.Offset{},
						EventTime: timestamppb.New(testEventTime),
						Keys:      []string{"test-key"},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "read_fn_err",
			fields: fields{
				readHandler: ReadFunc(func(ctx context.Context, request ReadRequest, messageCh chan<- model.Message) {
					msg := model.NewMessage([]byte(`test`), model.Offset{}, testEventTime)
					messageCh <- msg.WithKeys([]string{"test-key"})
					close(messageCh)
				}),
			},
			input: &sourcepb.ReadRequest{
				Request: &sourcepb.ReadRequest_Request{
					NumRecords:  1,
					TimeoutInMs: 1000,
				},
			},
			expected: []*sourcepb.DatumResponse{
				{
					Result: &sourcepb.DatumResponse_Result{
						Payload:   []byte(`test`),
						Offset:    &sourcepb.Offset{},
						EventTime: timestamppb.New(testEventTime),
						Keys:      []string{"test-key"},
					},
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				ReadHandler:    tt.fields.readHandler,
				AckHandler:     tt.fields.ackHandler,
				PendingHandler: tt.fields.pendingHandler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
			outputCh := make(chan sourcepb.DatumResponse)
			result := make([]*sourcepb.DatumResponse, 0)

			var readFnStream sourcepb.Source_ReadFnServer
			if tt.expectedErr {
				readFnStream = NewReadFnServerErrTest(ctx)
			} else {
				readFnStream = NewReadFnServerTest(ctx, outputCh)
			}

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					result = append(result, &msg)
				}
			}()

			err := fs.ReadFn(tt.input, readFnStream)
			close(outputCh)
			wg.Wait()

			if tt.expectedErr {
				assert.Error(t, err)
				// when the stream function returns an error, the message channel may or may not be closed.
				// so we skip asserting the result here.
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReadFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestService_AckFn(t *testing.T) {
	ackHandler := AckFunc(func(ctx context.Context, request AckRequest) {
		// Do nothing and return, to mimic a successful ack.
		return
	})
	fs := &Service{
		AckHandler: ackHandler,
	}
	ctx := context.Background()
	got, err := fs.AckFn(ctx, &sourcepb.AckRequest{
		Request: &sourcepb.AckRequest_Request{
			Offsets: []*sourcepb.Offset{
				{
					PartitionId: "0",
					Offset:      []byte("test"),
				},
			},
		},
	})
	assert.Equal(t, got, &sourcepb.AckResponse{
		Result: &sourcepb.AckResponse_Result{},
	})
	assert.NoError(t, err)
}

func TestService_PendingFn(t *testing.T) {
	pendingHandler := PendingFunc(func(ctx context.Context) uint64 {
		// Return 123, to mimic a pending func.
		return 123
	})
	fs := &Service{
		PendingHandler: pendingHandler,
	}
	ctx := context.Background()
	got, err := fs.PendingFn(ctx, &emptypb.Empty{})
	assert.Equal(t, got, &sourcepb.PendingResponse{
		Result: &sourcepb.PendingResponse_Result{
			Count: 123,
		},
	})
	assert.NoError(t, err)
}
