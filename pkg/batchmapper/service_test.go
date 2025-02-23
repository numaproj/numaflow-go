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

type BatchMapFnServerErrTest struct {
	ctx      context.Context
	inputCh  chan *mappb.MapRequest
	outputCh chan *mappb.MapResponse
	grpc.ServerStream
}

func NewBatchMapFnServerErrTest(
	ctx context.Context,
	inputCh chan *mappb.MapRequest,
	outputCh chan *mappb.MapResponse,

) *BatchMapFnServerErrTest {
	return &BatchMapFnServerErrTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *BatchMapFnServerErrTest) Recv() (*mappb.MapRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *BatchMapFnServerErrTest) Send(_ *mappb.MapResponse) error {
	return fmt.Errorf("send error")
}

func (u *BatchMapFnServerErrTest) Context() context.Context {
	return u.ctx
}

func TestService_BatchMapFn(t *testing.T) {
	tests := []struct {
		name        string
		handler     BatchMapper
		input       []*mappb.MapRequest
		expected    []*mappb.MapResponse
		expectedErr bool
	}{
		{
			name: "batch_map_stream_fn_forward_msg",
			handler: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for d := range datums {
					results := NewBatchResponse(d.ID())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*mappb.MapRequest{
				{
					Handshake: &mappb.Handshake{
						Sot: true,
					},
				},
				{
					Request: &mappb.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test1`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Id: "test1",
				},
				{
					Request: &mappb.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test2`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Id: "test2",
				},
				{
					Request: &mappb.MapRequest_Request{},
					Status: &mappb.TransmissionStatus{
						Eot: true,
					},
				},
			},
			expected: []*mappb.MapResponse{
				{
					Handshake: &mappb.Handshake{
						Sot: true,
					},
				},
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
				{
					Status: &mappb.TransmissionStatus{
						Eot: true,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "batch_map_stream_err",
			handler: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
				batchResponses := BatchResponsesBuilder()
				for d := range datums {
					results := NewBatchResponse(d.ID())
					results = results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
					batchResponses = batchResponses.Append(results)
				}
				return batchResponses
			}),
			input: []*mappb.MapRequest{
				{
					Handshake: &mappb.Handshake{
						Sot: true,
					},
				},
				{
					Request: &mappb.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test1`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Id: "test1",
				},
				{
					Request: &mappb.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test2`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
					Id: "test2",
				},
			},
			expected: []*mappb.MapResponse{
				{
					Handshake: &mappb.Handshake{
						Sot: true,
					},
				},
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
				{
					Status: &mappb.TransmissionStatus{
						Eot: true,
					},
				},
			},
			expectedErr: true,
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
			inputCh := make(chan *mappb.MapRequest, 3)
			outputCh := make(chan *mappb.MapResponse)
			result := make([]*mappb.MapResponse, 0)

			var udfBatchMapFnStream mappb.Map_MapFnServer
			if tt.expectedErr {
				udfBatchMapFnStream = NewBatchMapFnServerErrTest(ctx, inputCh, outputCh)
			} else {
				udfBatchMapFnStream = NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)
			}

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.MapFn(udfBatchMapFnStream)
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

			if tt.expectedErr {
				// assert err is not nil
				assert.NotNil(t, err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BatchMapFn() got = %v, want %v", result, tt.expected)
			}
		})
	}
}
