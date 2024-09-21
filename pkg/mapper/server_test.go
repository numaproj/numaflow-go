package mapper

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
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

	var mapHandler = MapperFunc(func(ctx context.Context, keys []string, d Datum) Messages {
		msg := d.Value()
		return MessagesBuilder().Append(NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
	})
	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := NewServer(mapHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}

// tests the case where the server is shutdown gracefully when a panic occurs in the map handler
func TestMapServer_GracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dir := t.TempDir()
	socketFile, _ := os.Create(dir + "/test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.Create(dir + "/numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var mapHandler = MapperFunc(func(ctx context.Context, keys []string, d Datum) Messages {
		msg := d.Value()
		if keys[0] == "key2" {
			time.Sleep(20 * time.Millisecond)
			panic("panic test")
		}
		time.Sleep(50 * time.Millisecond)
		return MessagesBuilder().Append(NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
	})

	done := make(chan struct{})
	go func() {
		err := NewServer(mapHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(socketFile.Name())).Start(ctx)
		assert.NoError(t, err)
		close(done)
	}()

	// wait for the server to start
	time.Sleep(10 * time.Millisecond)

	// create a client
	conn, err := grpc.NewClient(
		"unix://"+socketFile.Name(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := mappb.NewMapClient(conn)
	// send two map requests with key1 and key2 as keys simultaneously
	keys := []string{"key1", "key2"}
	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			req := &mappb.MapRequest{
				Keys: []string{key},
			}

			resp, err := client.MapFn(ctx, req)
			// since there is a panic in the map handler for key2, we should get an error
			// other requests should be successful
			if key == "key2" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		}(key)
	}

	wg.Wait()
	// wait for the server to shutdown gracefully because of the panic
	select {
	case <-ctx.Done():
		t.Fatal("server did not shutdown gracefully")
	case <-done:
	}
}
