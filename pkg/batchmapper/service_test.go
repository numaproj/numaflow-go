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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
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
					results := NewBatchResponse(d.Id())
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
							Keys:     []string{"client_test"},
							Value:    []byte(`test1`),
							Metadata: emptyResultMetadata(),
						},
					},
					Id: "test1",
				},
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:     []string{"client_test"},
							Value:    []byte(`test2`),
							Metadata: emptyResultMetadata(),
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
					results := NewBatchResponse(d.Id())
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
							Keys:     []string{"client_test"},
							Value:    []byte(`test1`),
							Metadata: emptyResultMetadata(),
						},
					},
					Id: "test1",
				},
				{
					Results: []*mappb.MapResponse_Result{
						{
							Keys:     []string{"client_test"},
							Value:    []byte(`test2`),
							Metadata: emptyResultMetadata(),
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

func TestService_BatchMapFn_Nack(t *testing.T) {
	delay := uint64(500)
	reason := "retry"
	fs := &Service{
		BatchMapper: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
			br := BatchResponsesBuilder()
			for d := range datums {
				r := NewBatchResponse(d.Id())
				r = r.Append(MessageToNack(&NackOptions{Delay: &delay, Reason: &reason}))
				br = br.Append(r)
			}
			return br
		}),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	inputCh := make(chan *mappb.MapRequest, 3)
	outputCh := make(chan *mappb.MapResponse)
	result := make([]*mappb.MapResponse, 0)
	stream := NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); _ = fs.MapFn(stream); close(outputCh) }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range outputCh {
			result = append(result, msg)
		}
	}()
	for _, val := range []*mappb.MapRequest{
		{Handshake: &mappb.Handshake{Sot: true}},
		{Request: &mappb.MapRequest_Request{Keys: []string{"k"}, Value: []byte("test")}, Id: "test1"},
		{Request: &mappb.MapRequest_Request{}, Status: &mappb.TransmissionStatus{Eot: true}},
	} {
		inputCh <- val
	}
	close(inputCh)
	wg.Wait()

	// result[0] = handshake, result[1] = the nacked message, result[2] = eot
	require.Len(t, result[1].Results, 1)
	assert.Equal(t, []string{NACK}, result[1].Results[0].Tags)
	assert.Equal(t, uint64(500), result[1].Results[0].GetNackOptions().GetDelay())
	assert.Equal(t, "retry", result[1].Results[0].GetNackOptions().GetReason())
}

func TestService_BatchMapFn_Fail(t *testing.T) {
	fs := &Service{
		BatchMapper: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
			br := BatchResponsesBuilder()
			for d := range datums {
				r := NewBatchResponse(d.Id())
				r = r.Append(MessageToFail())
				br = br.Append(r)
			}
			return br
		}),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	inputCh := make(chan *mappb.MapRequest, 3)
	outputCh := make(chan *mappb.MapResponse)
	result := make([]*mappb.MapResponse, 0)
	stream := NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); _ = fs.MapFn(stream); close(outputCh) }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range outputCh {
			result = append(result, msg)
		}
	}()
	for _, val := range []*mappb.MapRequest{
		{Handshake: &mappb.Handshake{Sot: true}},
		{Request: &mappb.MapRequest_Request{Keys: []string{"k"}, Value: []byte("test")}, Id: "test1"},
		{Request: &mappb.MapRequest_Request{}, Status: &mappb.TransmissionStatus{Eot: true}},
	} {
		inputCh <- val
	}
	close(inputCh)
	wg.Wait()

	// result[0] = handshake, result[1] = the failed message, result[2] = eot
	require.Len(t, result[1].Results, 1)
	assert.Equal(t, []string{FAIL}, result[1].Results[0].Tags)
}

func emptyResultMetadata() *common.Metadata {
	return &common.Metadata{
		SysMetadata:  map[string]*common.KeyValueGroup{},
		UserMetadata: map[string]*common.KeyValueGroup{},
	}
}

func TestService_BatchMapFn_Metadata(t *testing.T) {
	type seen struct {
		id          string
		traceparent string
		sysGroups   int
		userValue   string
	}
	var got []seen
	fs := &Service{
		BatchMapper: BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
			br := BatchResponsesBuilder()
			for d := range datums {
				sys := d.SystemMetadata()
				user := d.UserMetadata()
				entry := seen{id: d.Id()}
				if sys != nil {
					entry.traceparent = string(sys.Value("tracing_udf", "traceparent"))
					entry.sysGroups = len(sys.Groups())
				}
				if user != nil {
					entry.userValue = string(user.Value("app", "k"))
				}
				got = append(got, entry)
				msg := NewMessage(d.Value())
				if entry.userValue != "" {
					echoed := NewUserMetadata()
					echoed.AddKV("app", "k", []byte(entry.userValue))
					msg = msg.WithUserMetadata(echoed)
				}
				r := NewBatchResponse(d.Id())
				r = r.Append(msg)
				br = br.Append(r)
			}
			return br
		}),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	inputCh := make(chan *mappb.MapRequest, 4)
	outputCh := make(chan *mappb.MapResponse)
	result := make([]*mappb.MapResponse, 0)
	stream := NewBatchBatchMapStreamFnServerTest(ctx, inputCh, outputCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); _ = fs.MapFn(stream); close(outputCh) }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range outputCh {
			result = append(result, msg)
		}
	}()
	for _, val := range []*mappb.MapRequest{
		{Handshake: &mappb.Handshake{Sot: true}},
		{
			Request: &mappb.MapRequest_Request{
				Keys:  []string{"k"},
				Value: []byte("with-parent"),
				Metadata: &common.Metadata{
					SysMetadata: map[string]*common.KeyValueGroup{
						"tracing_udf": {KeyValue: map[string][]byte{"traceparent": []byte("00-trace-span-01")}},
					},
					UserMetadata: map[string]*common.KeyValueGroup{
						"app": {KeyValue: map[string][]byte{"k": []byte("v")}},
					},
				},
			},
			Id: "with-parent",
		},
		{Request: &mappb.MapRequest_Request{Keys: []string{"k"}, Value: []byte("no-parent")}, Id: "no-parent"},
		{Request: &mappb.MapRequest_Request{}, Status: &mappb.TransmissionStatus{Eot: true}},
	} {
		inputCh <- val
	}
	close(inputCh)
	wg.Wait()

	require.Equal(t, []seen{
		{id: "with-parent", traceparent: "00-trace-span-01", sysGroups: 1, userValue: "v"},
		{id: "no-parent", traceparent: "", sysGroups: 0},
	}, got)
	require.GreaterOrEqual(t, len(result), 3)
	assert.Equal(t, []byte("v"), result[1].Results[0].GetMetadata().GetUserMetadata()["app"].GetKeyValue()["k"])
	assert.Empty(t, result[1].Results[0].GetMetadata().GetSysMetadata())
	assert.Empty(t, result[2].Results[0].GetMetadata().GetUserMetadata())
	assert.Empty(t, result[2].Results[0].GetMetadata().GetSysMetadata())
}
