package server

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/client"
	"github.com/stretchr/testify/assert"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fields struct {
	mapHandler    functionsdk.MapHandler
	mapTHandler   functionsdk.MapTHandler
	reduceHandler functionsdk.ReduceHandler
}

func Test_server_map(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_map",
			fields: fields{
				mapHandler: functionsdk.MapFunc(func(ctx context.Context, key string, d functionsdk.Datum) functionsdk.Messages {
					msg := d.Value()
					return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(key+"_test", msg))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go New().RegisterMapper(tt.fields.mapHandler).Start(ctx, WithSockAddr(file.Name()))

			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("client_%d", i)
				// set the key in metadata for map function
				md := grpcmd.New(map[string]string{functionsdk.DatumKey: key})
				ctx = grpcmd.NewOutgoingContext(ctx, md)
				list, err := c.MapFn(ctx, &functionpb.Datum{
					Key:       key,
					Value:     []byte(`server_test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				})
				assert.NoError(t, err)
				for _, e := range list {
					assert.Equal(t, key+"_test", e.GetKey())
					assert.Equal(t, []byte(`server_test`), e.GetValue())
					assert.Nil(t, e.GetEventTime())
					assert.Nil(t, e.GetWatermark())
				}
			}
		})
	}
}

func Test_server_mapT(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_mapT",
			fields: fields{
				mapTHandler: functionsdk.MapTFunc(func(ctx context.Context, key string, d functionsdk.Datum) functionsdk.MessageTs {
					msg := d.Value()
					return functionsdk.MessageTsBuilder().Append(functionsdk.MessageTTo(time.Time{}, key+"_test", msg))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go New().RegisterMapperT(tt.fields.mapTHandler).Start(ctx, WithSockAddr(file.Name()))

			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("client_%d", i)
				// set the key in metadata for map function
				md := grpcmd.New(map[string]string{functionsdk.DatumKey: key})
				ctx = grpcmd.NewOutgoingContext(ctx, md)
				list, err := c.MapTFn(ctx, &functionpb.Datum{
					Key:       key,
					Value:     []byte(`server_test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				})
				assert.NoError(t, err)
				for _, e := range list {
					assert.Equal(t, key+"_test", e.GetKey())
					assert.Equal(t, []byte(`server_test`), e.GetValue())
					assert.Equal(t, &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})}, e.GetEventTime())
					assert.Nil(t, e.GetWatermark())
				}
			}
		})
	}
}

func Test_server_reduce(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()
	var testKey = "reduce_key"
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_reduce",
			fields: fields{
				reduceHandler: functionsdk.ReduceFunc(func(ctx context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
					// sum up values for the same key

					// in this test case, md is nil
					// intervalWindow := md.IntervalWindow()
					// _ = intervalWindow
					var resultKey = key
					var resultVal []byte
					var sum = 0
					// sum up the values
					for d := range reduceCh {
						val := d.Value()
						eventTime := d.EventTime()
						_ = eventTime
						watermark := d.Watermark()
						_ = watermark

						v, err := strconv.Atoi(string(val))
						if err != nil {
							continue
						}
						sum += v
					}
					resultVal = []byte(strconv.Itoa(sum))
					return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go New().RegisterReducer(tt.fields.reduceHandler).Start(ctx, WithSockAddr(file.Name()))

			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			var reduceDatumCh = make(chan *functionpb.Datum, 10)

			// the sum of the numbers from 0 to 9
			for i := 0; i < 10; i++ {
				reduceDatumCh <- &functionpb.Datum{
					Key:       testKey,
					Value:     []byte(strconv.Itoa(i)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				}
			}
			close(reduceDatumCh)

			// set the key in gPRC metadata for reduce function
			md := grpcmd.New(map[string]string{functionsdk.DatumKey: testKey, functionsdk.WinStartTime: "60000", functionsdk.WinEndTime: "120000"})
			ctx = grpcmd.NewOutgoingContext(ctx, md)
			list, err := c.ReduceFn(ctx, reduceDatumCh)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))
			assert.Equal(t, testKey, list[0].GetKey())
			assert.Equal(t, []byte(`45`), list[0].GetValue())
			assert.Nil(t, list[0].GetEventTime())
			assert.Nil(t, list[0].GetWatermark())
		})
	}
}
