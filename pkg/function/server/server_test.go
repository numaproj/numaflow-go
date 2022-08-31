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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_server_map(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.Remove(file.Name())
		assert.NoError(t, err)
	}()

	type fields struct {
		mapHandler    functionsdk.MapHandler
		reduceHandler functionsdk.ReduceHandler
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_map",
			fields: fields{
				mapHandler: functionsdk.MapFunc(func(ctx context.Context, key string, d functionsdk.Datum) (functionsdk.Messages, error) {
					msg := d.Value()
					return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(key+"_test", msg)), nil
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection

			go New().RegisterMapper(tt.fields.mapHandler).Start(context.Background(), WithSockAddr(file.Name()))

			var ctx = context.Background()
			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("client_%d", i)
				// set the key in metadata for reduce function
				md := metadata.New(map[string]string{functionsdk.DatumKey: key})
				ctx = metadata.NewOutgoingContext(ctx, md)
				list, err := c.MapFn(ctx, &functionpb.Datum{
					Key:   key,
					Value: []byte(`server_test`),
					EventTime: &functionpb.EventTime{
						EventTime: timestamppb.New(time.Unix(1661169600, 0)),
					},
					Watermark: &functionpb.Watermark{
						// TODO: need to update once we've finalized the datum data type
						Watermark: timestamppb.New(time.Time{}),
					},
				})
				assert.NoError(t, err)
				for _, e := range list {
					assert.Equal(t, key+"_test", e.Key)
					assert.Equal(t, []byte(`server_test`), e.Value)
				}
			}
		})
	}
}

func Test_server_reduce(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.Remove(file.Name())
		assert.NoError(t, err)
	}()

	var testKey = "reduce_key"

	type fields struct {
		mapHandler    functionsdk.MapHandler
		reduceHandler functionsdk.ReduceHandler
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_reduce",
			fields: fields{
				reduceHandler: functionsdk.ReduceFunc(func(ctx context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) (functionsdk.Messages, error) {
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
					return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal)), nil
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection

			go New().RegisterReducer(tt.fields.reduceHandler).Start(context.Background(), WithSockAddr(file.Name()))

			var ctx = context.Background()
			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			var reduceDatumCh = make(chan *functionpb.Datum, 10)

			for i := 0; i < 10; i++ {
				reduceDatumCh <- &functionpb.Datum{
					Key:   testKey,
					Value: []byte(strconv.Itoa(i)),
					EventTime: &functionpb.EventTime{
						EventTime: timestamppb.New(time.Unix(1661169600, 0)),
					},
					Watermark: &functionpb.Watermark{
						// TODO: need to update once we've finalized the datum data type
						Watermark: timestamppb.New(time.Time{}),
					},
				}
			}
			close(reduceDatumCh)

			// set the key in metadata for reduce function
			md := metadata.New(map[string]string{functionsdk.DatumKey: testKey})
			ctx = metadata.NewOutgoingContext(ctx, md)
			list, err := c.ReduceFn(ctx, reduceDatumCh)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))
			assert.Equal(t, testKey, list[0].GetKey())
			assert.Equal(t, []byte(`45`), list[0].GetValue())
		})
	}
}
