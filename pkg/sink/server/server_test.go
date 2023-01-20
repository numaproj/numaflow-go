package server

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/client"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_server_sink(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()

	type fields struct {
		sinkHandler sinksdk.SinkHandler
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_sink",
			fields: fields{
				sinkHandler: sinksdk.SinkFunc(func(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
					result := sinksdk.ResponsesBuilder()
					for d := range datumStreamCh {
						id := d.ID()
						if strings.Contains(string(d.Value()), "err") {
							result = result.Append(sinksdk.ResponseFailure(id, "mock sink message error"))
						} else {
							result = result.Append(sinksdk.ResponseOK(id))
						}

					}
					return result
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go New().RegisterSinker(tt.fields.sinkHandler).Start(ctx, WithSockAddr(file.Name()))

			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			testDatumList := []*sinkpb.Datum{
				{
					Id:        "test_id_0",
					Value:     []byte(`sink_message_success`),
					EventTime: &sinkpb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
					Watermark: &sinkpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Id:        "test_id_1",
					Value:     []byte(`sink_message_err`),
					EventTime: &sinkpb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
					Watermark: &sinkpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			}
			responseList, err := c.SinkFn(ctx, testDatumList)
			assert.NoError(t, err)
			expectedResponseList := []*sinkpb.Response{
				{
					Id:      "test_id_0",
					Success: true,
					ErrMsg:  "",
				},
				{
					Id:      "test_id_1",
					Success: false,
					ErrMsg:  "mock sink message error",
				},
			}
			assert.Equal(t, expectedResponseList, responseList)
		})
	}
}
