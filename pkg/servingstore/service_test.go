package servingstore

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	servingpb "github.com/numaproj/numaflow-go/pkg/apis/proto/serving/v1"
)

func newTestServer(t *testing.T, register func(server *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		_ = lis.Close()
	})

	server := grpc.NewServer()
	t.Cleanup(func() {
		server.Stop()
	})

	register(server)

	errChan := make(chan error, 1)
	go func() {
		// t.Fatal should only be called from the goroutine running the test
		if err := server.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		_ = conn.Close()
	})
	if err != nil {
		t.Fatalf("Creating new gRPC client connection: %v", err)
	}

	var grpcServerErr error
	select {
	case grpcServerErr = <-errChan:
	case <-time.After(500 * time.Millisecond):
		grpcServerErr = errors.New("gRPC server didn't start in 500ms")
	}
	if err != nil {
		t.Fatalf("Failed to start gRPC server: %v", grpcServerErr)
	}

	return conn
}

func TestService_Put(t *testing.T) {
	svc := &Service{
		ServingStore: newTestServingStore(false),
	}

	conn := newTestServer(t, func(server *grpc.Server) {
		servingpb.RegisterServingStoreServer(server, svc)
	})

	client := servingpb.NewServingStoreClient(conn)
	ctx := context.Background()

	req := &servingpb.PutRequest{
		Id: "test_id",
		Payloads: []*servingpb.Payload{
			{Origin: "test_origin", Value: []byte("test_value")},
		},
	}

	resp, err := client.Put(ctx, req)
	require.NoError(t, err, "Put request failed")
	assert.True(t, resp.Success, "Expected success response")
}

func TestService_Get(t *testing.T) {
	svc := &Service{
		ServingStore: &testServingStore{
			entries: map[string][]Payload{
				"test_id": {
					Payload{
						origin: "test_origin",
						value:  []byte("test_value"),
					},
				},
			},
			panicOnPut: false,
		},
	}

	conn := newTestServer(t, func(server *grpc.Server) {
		servingpb.RegisterServingStoreServer(server, svc)
	})

	client := servingpb.NewServingStoreClient(conn)
	ctx := context.Background()

	req := &servingpb.GetRequest{Id: "test_id"}

	resp, err := client.Get(ctx, req)
	require.NoError(t, err, "Get request failed")
	assert.Equal(t, "test_id", resp.Id, "Expected ID to match")
	assert.Len(t, resp.Payloads, 1, "Expected one payload")
	assert.Equal(t, "test_origin", resp.Payloads[0].Origin, "Expected origin to match")
	assert.Equal(t, []byte("test_value"), resp.Payloads[0].Value, "Expected value to match")
}

func TestService_Put_Panic(t *testing.T) {
	svc := &Service{
		ServingStore: newTestServingStore(true),
		shutdownCh:   make(chan<- struct{}, 1),
	}

	conn := newTestServer(t, func(server *grpc.Server) {
		servingpb.RegisterServingStoreServer(server, svc)
	})

	client := servingpb.NewServingStoreClient(conn)
	ctx := context.Background()

	req := &servingpb.PutRequest{
		Id: "test_id",
		Payloads: []*servingpb.Payload{
			{Origin: "test_origin", Value: []byte("test_value")},
		},
	}

	_, err := client.Put(ctx, req)
	require.Error(t, err, "Expected error while sending Put request")
	gotStatus, _ := status.FromError(err)
	gotMessage := gotStatus.Message()
	expectedStatus := status.Convert(status.Errorf(codes.Internal, "%s: %v", errServingStorePanic, "put failed"))
	expectedMessage := expectedStatus.Message()
	assert.Equal(t, expectedStatus.Code(), gotStatus.Code(), "Expected error codes to be equal")
	assert.True(t, strings.HasPrefix(gotMessage, expectedMessage), "Expected error message to start with the expected message")
}

type testServingStore struct {
	entries    map[string][]Payload
	panicOnPut bool
}

func newTestServingStore(panicOnPut bool) *testServingStore {
	return &testServingStore{
		entries:    make(map[string][]Payload),
		panicOnPut: panicOnPut,
	}
}

func (t *testServingStore) Put(ctx context.Context, put PutDatum) {
	if t.panicOnPut {
		panic("put failed")
	}
	t.entries[put.ID()] = put.Payloads()
}

func (t *testServingStore) Get(ctx context.Context, get GetDatum) StoredResult {
	return StoredResult{
		id:       get.ID(),
		payloads: t.entries[get.ID()],
	}
}
