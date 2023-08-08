package server

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
)

func TestServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var mapHandler = functionsdk.MapFunc(func(ctx context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
		msg := d.Value()
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := New().RegisterMapper(mapHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
	assert.NoError(t, err)
}

func TestServer_RegisterMapper(t *testing.T) {
	var mapHandler = functionsdk.MapFunc(func(ctx context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
		msg := d.Value()
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
	})

	serv := New()
	assert.Nil(t, serv.svc.Mapper)
	serv.RegisterMapper(mapHandler)
	assert.NotNil(t, serv.svc.Mapper)
}

func TestServer_RegisterMapperStream(t *testing.T) {
	mapStreamHandler := functionsdk.MapStreamFunc(func(ctx context.Context, keys []string, d functionsdk.Datum, messageCh chan<- functionsdk.Message) {
		msg := d.Value()
		messageCh <- functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
		close(messageCh)
	})
	serv := New()
	assert.Nil(t, serv.svc.MapperStream)
	serv.RegisterMapperStream(mapStreamHandler)
	assert.NotNil(t, serv.svc.MapperStream)
}

func TestServer_RegisterMapperT(t *testing.T) {
	mapTHandler := functionsdk.MapTFunc(func(ctx context.Context, keys []string, d functionsdk.Datum) functionsdk.MessageTs {
		msg := d.Value()
		return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(msg, time.Time{}).WithKeys([]string{keys[0] + "_test"}))
	})
	serv := New()
	assert.Nil(t, serv.svc.MapperT)
	serv.RegisterMapperT(mapTHandler)
	assert.NotNil(t, serv.svc.MapperT)
}

func TestServer_RegisterReducer(t *testing.T) {
	reduceHandler := functionsdk.ReduceFunc(func(ctx context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
		var resultKey = keys
		var resultVal []byte
		var sum = 0
		resultVal = []byte(strconv.Itoa(sum))
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKey))
	})
	serv := New()
	assert.Nil(t, serv.svc.Reducer)
	serv.RegisterReducer(reduceHandler)
	assert.NotNil(t, serv.svc.Reducer)
}
