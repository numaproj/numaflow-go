package mapper

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

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

type BatchMapStreamFnServerTest struct {
	ctx      context.Context
	outputCh chan *mappb.MapResponse
	inputCh  chan *mappb.MapRequest
	grpc.ServerStream
}

func NewBatchBatchMapStreamFnServerTest(
	ctx context.Context,
	inputCh chan *mappb.MapRequest,
	outputCh chan *mappb.MapResponse,
) *BatchMapStreamFnServerTest {
	return &BatchMapStreamFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *BatchMapStreamFnServerTest) Recv() (*mappb.MapRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *BatchMapStreamFnServerTest) Send(d *mappb.MapResponse) error {
	u.outputCh <- d
	return nil
}

func (u *BatchMapStreamFnServerTest) Context() context.Context {
	return u.ctx
}

type MapStreamFnServerErrTest struct {
	ctx      context.Context
	inputCh  chan *mappb.MapRequest
	outputCh chan *mappb.MapResponse
	grpc.ServerStream
}

func NewMapStreamFnServerErrTest(
	ctx context.Context,
	inputCh chan *mappb.MapRequest,
	outputCh chan *mappb.MapResponse,

) *MapStreamFnServerErrTest {
	return &MapStreamFnServerErrTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *MapStreamFnServerErrTest) Recv() (*mappb.MapRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *MapStreamFnServerErrTest) Send(_ *mappb.MapResponse) error {
	return fmt.Errorf("send error")
}

func (u *MapStreamFnServerErrTest) Context() context.Context {
	return u.ctx
}

func TestService_MapFnStream(t *testing.T) {
	tests := []struct {
		name        string
		handler     BatchMapper
		input       []*mappb.MapRequest
		expected    []*mappb.MapResponse
		expectedErr bool
		streamErr   bool
	}{
		{
			name: "map_stream_fn_forward_msg",
			handler: BatchMapperFunc(func(ctx context.Context, datums []Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for _, d := range datums {
					results := NewBatchResponse(d.Id())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*mappb.MapRequest{{
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
			expected: []*mappb.MapResponse{
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test1`),
						},
					},
					Id: "test1",
				},
				{
					Results: []*mappb.MapResponse_Result{
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
			name: "batch_map_mismatch_output_len",
			handler: BatchMapperFunc(func(ctx context.Context, datums []Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				return batchResponses
			}),
			input: []*mappb.MapRequest{{
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
			expected: []*mappb.MapResponse{
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test1`),
						},
					},
					Id: "test1",
				},
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test2`),
						},
					},
					Id: "test2",
				},
			},
			expectedErr: true,
		},
		{
			name: "batch_map_stream_err",
			handler: BatchMapperFunc(func(ctx context.Context, datums []Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for _, d := range datums {
					results := NewBatchResponse(d.Id())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*mappb.MapRequest{{
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
			expected: []*mappb.MapResponse{
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:  []string{"client_test"},
							Value: []byte(`test1`),
						},
					},
					Id: "test1",
				},
				{
					Results: []*mappb.MapResponse_Result{
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
			ctx := context.Background()
			inputCh := make(chan *mappb.MapRequest)
			outputCh := make(chan *mappb.MapResponse)
			result := make([]*mappb.MapResponse, 0)

			var udfMapStreamFnStream mappb.Map_MapStreamFnServer
			if tt.streamErr {
				udfMapStreamFnStream = NewMapStreamFnServerErrTest(ctx, inputCh, outputCh)
			} else {
				udfMapStreamFnStream = NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)
			}

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.MapStreamFn(udfMapStreamFnStream)
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
				assert.True(t, tt.expectedErr, "MapStreamFn() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapStreamFn() got = %v, want %v", result, tt.expected)
			}

		})
	}
}
