package server

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
)

func TestMapServer_Start(t *testing.T) {
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
	err := NewMapServer(mapHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}

func TestMapStreamServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var mapStreamHandler = functionsdk.MapStreamFunc(func(ctx context.Context, keys []string, datum functionsdk.Datum, messageCh chan<- functionsdk.Message) {
		msg := datum.Value()
		messageCh <- functionsdk.NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
		close(messageCh)
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := NewMapStreamServer(ctx, mapStreamHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}

func TestReduceServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var reduceHandler = functionsdk.ReduceFunc(func(ctx context.Context, keys []string, rch <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
		sum := 0
		for val := range rch {
			msgVal, _ := strconv.Atoi(string(val.Value()))
			sum += msgVal
		}
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := NewReduceServer(reduceHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
