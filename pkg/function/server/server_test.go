package server

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow/pkg/func-client/client"
)

func waitUntilReady(ctx context.Context, c client.Client, t *testing.T) {
	var in = &emptypb.Empty{}
	isReady, _ := c.IsReady(ctx, in)
	for !isReady {
		select {
		case <-ctx.Done():
			assert.Fail(t, ctx.Err().Error())
			return
		default:
			time.Sleep(100 * time.Millisecond)
			isReady, _ = c.IsReady(ctx, in)
		}
	}
}

type fields struct {
	mapHandler       functionsdk.MapHandler
	mapTHandler      functionsdk.MapTHandler
	reduceHandler    functionsdk.ReduceHandler
	mapStreamHandler functionsdk.MapStreamFunc
}

func Test_server_map(t *testing.T) {
	socketFile, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(socketFile.Name())
		assert.NoError(t, err)
	}()

	serverInfoFile, err := os.CreateTemp("/tmp", "numaflow-test-info")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(serverInfoFile.Name())
		assert.NoError(t, err)
	}()

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_map",
			fields: fields{
				mapHandler: functionsdk.MapFunc(func(ctx context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
					msg := d.Value()
					return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			go New().RegisterMapper(tt.fields.mapHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
			c, err := client.New(client.WithUdsSockAddr(socketFile.Name()), client.WithServerInfoFilePath(serverInfoFile.Name()))
			waitUntilReady(ctx, c, t)
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()

			for i := 0; i < 10; i++ {
				keys := []string{fmt.Sprintf("client_%d", i)}
				list, err := c.MapFn(ctx, &functionpb.DatumRequest{
					Keys:      keys,
					Value:     []byte(`server_test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				})
				assert.NoError(t, err)
				for _, e := range list {
					assert.Equal(t, []string{keys[0] + "_test"}, e.GetKeys())
					assert.Equal(t, []byte(`server_test`), e.GetValue())
					assert.Nil(t, e.GetEventTime())
					assert.Nil(t, e.GetWatermark())
				}
			}
		})
	}
}

func Test_server_map_stream(t *testing.T) {
	socketFile, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(socketFile.Name())
		assert.NoError(t, err)
	}()

	serverInfoFile, err := os.CreateTemp("/tmp", "numaflow-test-info")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(serverInfoFile.Name())
		assert.NoError(t, err)
	}()

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_map",
			fields: fields{
				mapStreamHandler: functionsdk.MapStreamFunc(func(ctx context.Context, keys []string, d functionsdk.Datum, messageCh chan<- functionsdk.Message) {
					msg := d.Value()
					messageCh <- functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
					close(messageCh)
				}),
			},
		},
		{
			name: "server_map_without_close_stream",
			fields: fields{
				mapStreamHandler: functionsdk.MapStreamFunc(func(ctx context.Context, keys []string, d functionsdk.Datum, messageCh chan<- functionsdk.Message) {
					msg := d.Value()
					messageCh <- functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			go New().RegisterMapperStream(tt.fields.mapStreamHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))

			c, err := client.New(client.WithUdsSockAddr(socketFile.Name()), client.WithServerInfoFilePath(serverInfoFile.Name()))
			waitUntilReady(ctx, c, t)
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()

			for i := 0; i < 10; i++ {
				keys := []string{fmt.Sprintf("client_%d", i)}

				datumCh := make(chan *functionpb.DatumResponse)
				go func() {
					err := c.MapStreamFn(ctx, &functionpb.DatumRequest{
						Keys:      keys,
						Value:     []byte(`server_test`),
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
						Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
					}, datumCh)
					assert.NoError(t, err)
				}()

				for msg := range datumCh {
					assert.Equal(t, []string{keys[0] + "_test"}, msg.GetKeys())
					assert.Equal(t, []byte(`server_test`), msg.GetValue())
					assert.Nil(t, msg.GetEventTime())
					assert.Nil(t, msg.GetWatermark())
				}
			}
		})
	}
}

func Test_server_mapT(t *testing.T) {
	socketFile, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(socketFile.Name())
		assert.NoError(t, err)
	}()

	serverInfoFile, err := os.CreateTemp("/tmp", "numaflow-test-info")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(serverInfoFile.Name())
		assert.NoError(t, err)
	}()

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_mapT",
			fields: fields{
				mapTHandler: functionsdk.MapTFunc(func(ctx context.Context, keys []string, d functionsdk.Datum) functionsdk.MessageTs {
					msg := d.Value()
					return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(msg, time.Time{}).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			go New().RegisterMapperT(tt.fields.mapTHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
			c, err := client.New(client.WithUdsSockAddr(socketFile.Name()), client.WithServerInfoFilePath(serverInfoFile.Name()))
			waitUntilReady(ctx, c, t)
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			for i := 0; i < 10; i++ {
				keys := []string{fmt.Sprintf("client_%d", i)}
				list, err := c.MapTFn(ctx, &functionpb.DatumRequest{
					Keys:      keys,
					Value:     []byte(`server_test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				})
				assert.NoError(t, err)
				for _, e := range list {
					assert.Equal(t, []string{keys[0] + "_test"}, e.GetKeys())
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

	serverInfoFile, err := os.CreateTemp("/tmp", "numaflow-test-info")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(serverInfoFile.Name())
		assert.NoError(t, err)
	}()

	var testKeys = []string{"reduce_key"}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_reduce",
			fields: fields{
				reduceHandler: functionsdk.ReduceFunc(func(ctx context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
					// sum up values for the same keys

					// in this test case, md is nil
					// intervalWindow := md.IntervalWindow()
					// _ = intervalWindow
					var resultKey = keys
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
					return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKey))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			go New().RegisterReducer(tt.fields.reduceHandler).Start(ctx, WithSockAddr(file.Name()), WithServerInfoFilePath(serverInfoFile.Name()))

			c, err := client.New(client.WithUdsSockAddr(file.Name()), client.WithServerInfoFilePath(serverInfoFile.Name()))
			waitUntilReady(ctx, c, t)
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			var reduceDatumCh = make(chan *functionpb.DatumRequest, 10)

			// the sum of the numbers from 0 to 9
			for i := 0; i < 10; i++ {
				reduceDatumCh <- &functionpb.DatumRequest{
					Keys:      testKeys,
					Value:     []byte(strconv.Itoa(i)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				}
			}
			close(reduceDatumCh)

			dList := &functionpb.DatumResponseList{}
			md := grpcmd.New(map[string]string{functionsdk.WinStartTime: "60000", functionsdk.WinEndTime: "120000"})
			ctx = grpcmd.NewOutgoingContext(ctx, md)
			resultDatumList, err := c.ReduceFn(ctx, reduceDatumCh)
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				dList.Elements = append(dList.Elements, resultDatumList...)
			}()

			wg.Wait()
			assert.NoError(t, err)
			assert.Equal(t, 1, len(dList.Elements))
			assert.Equal(t, testKeys, dList.Elements[0].GetKeys())
			assert.Equal(t, []byte(`45`), dList.Elements[0].GetValue())
			assert.Nil(t, dList.Elements[0].GetEventTime())
			assert.Nil(t, dList.Elements[0].GetWatermark())
		})
	}
}
