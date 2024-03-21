package sourcer

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
)

var testEventTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var testKey = "test-key"
var testPendingNumber int64 = 123
var testPartitions = []int32{1, 3, 5}

type TestSource struct{}

func (ts TestSource) Read(_ context.Context, _ ReadRequest, messageCh chan<- Message) {
	msg := NewMessage([]byte(`test`), Offset{}, testEventTime).WithHeaders(map[string]string{"x-txn-id": "test-txn-id"})
	messageCh <- msg.WithKeys([]string{testKey})
}

func (ts TestSource) Ack(_ context.Context, _ AckRequest) {
	// Do nothing and return, to mimic a successful ack.
	return
}

func (ts TestSource) Pending(_ context.Context) int64 {
	return testPendingNumber
}

func (ts TestSource) Partitions(_ context.Context) []int32 {
	return testPartitions
}

func TestService_IsReady(t *testing.T) {
	fs := &Service{
		Source: nil,
	}
	got, err := fs.IsReady(nil, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, got, &sourcepb.ReadyResponse{
		Ready: true,
	})
}

type ReadFnServerTest struct {
	ctx      context.Context
	outputCh chan sourcepb.ReadResponse
	grpc.ServerStream
}

func NewReadFnServerTest(
	ctx context.Context,
	outputCh chan sourcepb.ReadResponse,
) *ReadFnServerTest {
	return &ReadFnServerTest{
		ctx:      ctx,
		outputCh: outputCh,
	}
}

func (t *ReadFnServerTest) Send(d *sourcepb.ReadResponse) error {
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

func (te *ReadFnServerErrTest) Send(_ *sourcepb.ReadResponse) error {
	return fmt.Errorf("send error")
}

func (te *ReadFnServerErrTest) Context() context.Context {
	return te.ctx
}

func TestService_ReadFn(t *testing.T) {
	tests := []struct {
		name        string
		input       *sourcepb.ReadRequest
		expected    []*sourcepb.ReadResponse
		expectedErr bool
	}{
		{
			name: "read_fn_read_msg",
			input: &sourcepb.ReadRequest{
				Request: &sourcepb.ReadRequest_Request{
					NumRecords:  1,
					TimeoutInMs: 1000,
				},
			},
			expected: []*sourcepb.ReadResponse{
				{
					Result: &sourcepb.ReadResponse_Result{
						Payload:   []byte(`test`),
						Offset:    &sourcepb.Offset{},
						EventTime: timestamppb.New(testEventTime),
						Keys:      []string{testKey},
						Headers:   map[string]string{"x-txn-id": "test-txn-id"},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "read_fn_err",
			input: &sourcepb.ReadRequest{
				Request: &sourcepb.ReadRequest_Request{
					NumRecords:  1,
					TimeoutInMs: 1000,
				},
			},
			expected: []*sourcepb.ReadResponse{
				{
					Result: &sourcepb.ReadResponse_Result{
						Payload:   []byte(`test`),
						Offset:    &sourcepb.Offset{},
						EventTime: timestamppb.New(testEventTime),
						Keys:      []string{testKey},
						Headers:   map[string]string{"x-txn-id": "test-txn-id"},
					},
				},
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{Source: TestSource{}}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
			outputCh := make(chan sourcepb.ReadResponse)
			result := make([]*sourcepb.ReadResponse, 0)

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
	fs := &Service{Source: TestSource{}}
	ctx := context.Background()
	got, err := fs.AckFn(ctx, &sourcepb.AckRequest{
		Request: &sourcepb.AckRequest_Request{
			Offsets: []*sourcepb.Offset{
				{
					PartitionId: 0,
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
	fs := &Service{Source: TestSource{}}
	ctx := context.Background()
	got, err := fs.PendingFn(ctx, &emptypb.Empty{})
	assert.Equal(t, got, &sourcepb.PendingResponse{
		Result: &sourcepb.PendingResponse_Result{
			Count: testPendingNumber,
		},
	})
	assert.NoError(t, err)
}

func TestService_PartitionsFn(t *testing.T) {
	fs := &Service{Source: TestSource{}}
	ctx := context.Background()
	got, err := fs.PartitionsFn(ctx, &emptypb.Empty{})
	assert.EqualValues(t, got, &sourcepb.PartitionsResponse{
		Result: &sourcepb.PartitionsResponse_Result{
			Partitions: testPartitions,
		},
	})
	assert.NoError(t, err)
}
