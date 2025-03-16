package accumulator

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/protobuf/types/known/emptypb"

	accumulatorpb "github.com/numaproj/numaflow-go/pkg/apis/proto/accumulator/v1"
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

func TestService_IsReady(t *testing.T) {
	svc := &Service{}
	conn := newTestServer(t, func(server *grpc.Server) {
		accumulatorpb.RegisterAccumulatorServer(server, svc)
	})

	client := accumulatorpb.NewAccumulatorClient(conn)
	resp, err := client.IsReady(context.Background(), &emptypb.Empty{})
	require.NoError(t, err, "Calling IsReady")
	assert.True(t, resp.Ready, "Expected IsReady to return true")
}

func TestService_AccumulateFn(t *testing.T) {
	svc := &Service{
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Datum) {
			for datum := range input {
				output <- datum
			}
		}),
		shutdownCh: make(chan struct{}, 1),
	}
	conn := newTestServer(t, func(server *grpc.Server) {
		accumulatorpb.RegisterAccumulatorServer(server, svc)
	})

	client := accumulatorpb.NewAccumulatorClient(conn)
	stream, err := client.AccumulateFn(context.Background())
	require.NoError(t, err, "Creating stream")

	// Send a request
	req := &accumulatorpb.AccumulatorRequest{
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event: accumulatorpb.AccumulatorRequest_WindowOperation_APPEND,
		},
	}
	err = stream.Send(req)
	require.NoError(t, err, "Sending request")

	// Receive a response
	resp, err := stream.Recv()
	require.NoError(t, err, "Receiving response")
	assert.NotNil(t, resp, "Expected a non-nil response")

	// Close the stream
	err = stream.CloseSend()
	require.NoError(t, err, "Closing the send direction of the stream")
}

func TestService_AccumulateFn_Panic(t *testing.T) {
	panicMssg := "accumulate failed"
	svc := &Service{
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Datum) {
			panic(panicMssg)
		}),
		shutdownCh: make(chan struct{}, 1),
	}
	conn := newTestServer(t, func(server *grpc.Server) {
		accumulatorpb.RegisterAccumulatorServer(server, svc)
	})

	client := accumulatorpb.NewAccumulatorClient(conn)
	stream, err := client.AccumulateFn(context.Background())
	require.NoError(t, err, "Creating stream")

	// Send a request
	req := &accumulatorpb.AccumulatorRequest{
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event: accumulatorpb.AccumulatorRequest_WindowOperation_APPEND,
		},
	}
	err = stream.Send(req)
	require.NoError(t, err, "Sending request")

	require.NoError(t, err, "Closing the send direction of the stream")

	// Receive a response
	_, err = stream.Recv()
	require.Error(t, err, "Expected error while receiving response")
	gotStatus, _ := status.FromError(err)
	gotMessage := gotStatus.Message()
	expectedStatus := status.Convert(status.Errorf(codes.Internal, "%s: %v", errAccumulatorPanic, panicMssg))
	expectedMessage := expectedStatus.Message()
	require.Equal(t, expectedStatus.Code(), gotStatus.Code(), "Expected error codes to be equal")
	require.True(t, gotMessage == expectedMessage, "Expected error message to match the expected message")
}
