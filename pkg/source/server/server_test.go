package server

import (
	"context"
	"fmt"
	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"os"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/source/client"
	"github.com/stretchr/testify/assert"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_server_transform(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()

	type fields struct {
		transformHandler sourcesdk.TransformHandler
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "server_transform",
			fields: fields{
				transformHandler: sourcesdk.TransformFunc(func(ctx context.Context, key string, d sourcesdk.Datum) sourcesdk.Messages {
					msg := d.Value()
					return sourcesdk.MessagesBuilder().Append(sourcesdk.MessageTo(time.Time{}, key+"_test", msg))
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// note: using actual UDS connection
			go New().RegisterTransformer(tt.fields.transformHandler).Start(context.Background(), WithSockAddr(file.Name()))

			var ctx = context.Background()
			c, err := client.New(client.WithSockAddr(file.Name()))
			assert.NoError(t, err)
			defer func() {
				err = c.CloseConn(ctx)
				assert.NoError(t, err)
			}()
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("client_%d", i)
				md := grpcmd.New(map[string]string{sourcesdk.DatumKey: key})
				ctx = grpcmd.NewOutgoingContext(ctx, md)
				list, err := c.TransformFn(ctx, &sourcepb.Datum{
					Key:       key,
					Value:     []byte(`server_test`),
					EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &sourcepb.Watermark{Watermark: timestamppb.New(time.Time{})},
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(list))
				for _, e := range list {
					assert.Equal(t, key+"_test", e.GetKey())
					assert.Equal(t, []byte(`server_test`), e.GetValue())
					assert.Equal(t, &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})}, e.GetEventTime())
					assert.Nil(t, e.GetWatermark())
				}
			}
		})
	}
}
