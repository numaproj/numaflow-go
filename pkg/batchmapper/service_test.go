package batchmapper

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	batchmappb "github.com/numaproj/numaflow-go/pkg/apis/proto/batchmap/v1"
)

type BatchMapStreamFnServerTest struct {
	ctx      context.Context
	outputCh chan *batchmappb.BatchMapResponse
	inputCh  chan *batchmappb.BatchMapRequest
	grpc.ServerStream
}

func NewBatchBatchMapStreamFnServerTest(
	ctx context.Context,
	inputCh chan *batchmappb.BatchMapRequest,
	outputCh chan *batchmappb.BatchMapResponse,
) *BatchMapStreamFnServerTest {
	return &BatchMapStreamFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *BatchMapStreamFnServerTest) Recv() (*batchmappb.BatchMapRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *BatchMapStreamFnServerTest) Send(d *batchmappb.BatchMapResponse) error {
	u.outputCh <- d
	return nil
}

func (u *BatchMapStreamFnServerTest) Context() context.Context {
	return u.ctx
}

type BatchMapFnServerErrTest struct {
	ctx      context.Context
	inputCh  chan *batchmappb.BatchMapRequest
	outputCh chan *batchmappb.BatchMapResponse
	grpc.ServerStream
}

func NewBatchMapFnServerErrTest(
	ctx context.Context,
	inputCh chan *batchmappb.BatchMapRequest,
	outputCh chan *batchmappb.BatchMapResponse,

) *BatchMapFnServerErrTest {
	return &BatchMapFnServerErrTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *BatchMapFnServerErrTest) Recv() (*batchmappb.BatchMapRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *BatchMapFnServerErrTest) Send(_ *batchmappb.BatchMapResponse) error {
	return fmt.Errorf("send error")
}

func (u *BatchMapFnServerErrTest) Context() context.Context {
	return u.ctx
}

func TestService_BatchMapFn(t *testing.T) {
	tests := []struct {
		name        string
		handler     BatchMapper
		input       []*batchmappb.BatchMapRequest
		expected    []*batchmappb.BatchMapResponse
		expectedErr bool
		streamErr   bool
	}{
		{
			name: "batch_map_stream_fn_forward_msg",
			handler: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for d := range datums {
					results := NewBatchResponse(d.Id())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*batchmappb.BatchMapRequest{{
				Keys:      []string{"client"},
				Value:     []byte(`test1`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
				Id:        "test1",
			}, {
				Keys:      []string{"client"},
				Value:     []byte(`test2`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
				Id:        "test2",
			}},
			expected: []*batchmappb.BatchMapResponse{
				{
					Results: []*batchmappb.BatchMapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test1`),
						},
					},
					Id: "test1",
				},
				{
					Results: []*batchmappb.BatchMapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test2`),
						},
					},
					Id: "test2",
				},
			},
			expectedErr: false,
		},
		{
			name: "batch_map_stream_err",
			handler: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for d := range datums {
					results := NewBatchResponse(d.Id())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*batchmappb.BatchMapRequest{{
				Keys:      []string{"client"},
				Value:     []byte(`test1`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
				Id:        "test1",
			}, {
				Keys:      []string{"client"},
				Value:     []byte(`test2`),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
				Id:        "test2",
			}},
			expected: []*batchmappb.BatchMapResponse{
				{
					Results: []*batchmappb.BatchMapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test1`),
						},
					},
					Id: "test1",
				},
				{
					Results: []*batchmappb.BatchMapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test2`),
						},
					},
					Id: "test2",
				},
			},
			expectedErr: true,
			streamErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				BatchMapper: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			inputCh := make(chan *batchmappb.BatchMapRequest)
			outputCh := make(chan *batchmappb.BatchMapResponse)
			result := make([]*batchmappb.BatchMapResponse, 0)

			var udfBatchMapFnStream batchmappb.BatchMap_BatchMapFnServer
			if tt.streamErr {
				udfBatchMapFnStream = NewBatchMapFnServerErrTest(ctx, inputCh, outputCh)
			} else {
				udfBatchMapFnStream = NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)
			}

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.BatchMapFn(udfBatchMapFnStream)
				close(outputCh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					result = append(result, msg)
				}
			}()

			for _, val := range tt.input {
				inputCh <- val
			}
			close(inputCh)
			wg.Wait()

			if err != nil {
				assert.True(t, tt.expectedErr, "BatchMapFn() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BatchMapFn() got = %v, want %v", result, tt.expected)
			}

		})
	}
}
