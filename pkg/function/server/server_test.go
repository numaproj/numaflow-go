package server

import (
	"context"
	"fmt"
	"os"
	"testing"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/client"
	"github.com/stretchr/testify/assert"
)

func Test_server_Start(t *testing.T) {
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
			name: "server_start",
			fields: fields{
				mapHandler: functionsdk.MapFunc(func(ctx context.Context, d functionsdk.Datum) (functionsdk.Messages, error) {
					key := d.Key()
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
				list, err := c.MapFn(ctx, &functionpb.Datum{
					Key:   key,
					Value: []byte(`server_test`),
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
