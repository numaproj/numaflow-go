package accumulator

import (
	"context"
	"errors"
	"io"
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
	"google.golang.org/protobuf/types/known/timestamppb"

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
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Message) {
			for datum := range input {
				output <- MessageFromDatum(datum)
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
	shutdownCh := make(chan struct{}, 1)
	svc := &Service{
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Message) {
			panic(panicMssg)
		}),
		shutdownCh: shutdownCh,
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

	// Give a small delay to ensure error propagation completes
	time.Sleep(100 * time.Millisecond)

	err = stream.CloseSend()
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

	// Verify that shutdown channel was triggered
	select {
	case <-shutdownCh:
		// Expected - shutdown was triggered
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected shutdown channel to be triggered, but it wasn't")
	}
}

func TestService_AccumulateFn_EOFEchoesCloseWindow(t *testing.T) {
	svc := &Service{
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Message) {
			for datum := range input {
				output <- MessageFromDatum(datum)
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

	keys := []string{"k1"}
	openWindow := &accumulatorpb.KeyedWindow{
		Start: timestamppb.New(time.UnixMilli(0)),
		End:   timestamppb.New(time.UnixMilli(60000)),
		Slot:  "slot-0",
		Keys:  keys,
	}

	// OPEN
	require.NoError(t, stream.Send(&accumulatorpb.AccumulatorRequest{
		Payload: &accumulatorpb.Payload{
			Keys:      keys,
			Value:     []byte("v1"),
			EventTime: timestamppb.New(time.UnixMilli(1000)),
			Watermark: timestamppb.New(time.UnixMilli(1000)),
		},
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event:       accumulatorpb.AccumulatorRequest_WindowOperation_OPEN,
			KeyedWindow: openWindow,
		},
	}), "Sending OPEN")

	// APPEND
	require.NoError(t, stream.Send(&accumulatorpb.AccumulatorRequest{
		Payload: &accumulatorpb.Payload{
			Keys:      keys,
			Value:     []byte("v2"),
			EventTime: timestamppb.New(time.UnixMilli(2000)),
			Watermark: timestamppb.New(time.UnixMilli(2000)),
		},
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event:       accumulatorpb.AccumulatorRequest_WindowOperation_APPEND,
			KeyedWindow: openWindow,
		},
	}), "Sending APPEND")

	// CLOSE with a distinct window that must be echoed in the EOF response.
	closeWindow := &accumulatorpb.KeyedWindow{
		Start: timestamppb.New(time.UnixMilli(1000000)),
		End:   timestamppb.New(time.UnixMilli(2000000)),
		Slot:  "slot-7",
		Keys:  keys,
	}
	require.NoError(t, stream.Send(&accumulatorpb.AccumulatorRequest{
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event:       accumulatorpb.AccumulatorRequest_WindowOperation_CLOSE,
			KeyedWindow: closeWindow,
		},
	}), "Sending CLOSE")

	// Read responses on the still-open stream until we get the EOF response, mirroring how
	// numaflow-core reads the per-window EOF while the bidi stream stays open. Assert it
	// echoes the CLOSE window.
	var eof *accumulatorpb.AccumulatorResponse
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "Receiving response")
		if resp.GetEOF() {
			eof = resp
			break
		}
	}

	require.NoError(t, stream.CloseSend(), "CloseSend")

	require.NotNil(t, eof, "Expected an EOF response")
	w := eof.GetWindow()
	require.NotNil(t, w, "EOF response must carry a window")
	assert.Equal(t, closeWindow.GetStart().GetSeconds(), w.GetStart().GetSeconds(), "EOF window start")
	assert.Equal(t, closeWindow.GetEnd().GetSeconds(), w.GetEnd().GetSeconds(), "EOF window end")
	assert.Equal(t, "slot-7", w.GetSlot(), "EOF window slot")
	assert.Equal(t, keys, w.GetKeys(), "EOF window keys")
}

// TestService_AccumulateFn_EOFEchoesCloseWindow_NoOutput exercises the WAL-leak scenario:
// an accumulator that emits nothing must still echo the CLOSE window's end in its EOF
// response (rather than the never-advanced latest watermark) so core can GC the WAL.
func TestService_AccumulateFn_EOFEchoesCloseWindow_NoOutput(t *testing.T) {
	svc := &Service{
		AccumulatorCreator: SimpleCreatorWithAccumulateFn(func(ctx context.Context, input <-chan Datum, output chan<- Message) {
			for range input {
				// intentionally emit nothing
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

	keys := []string{"k1"}
	openWindow := &accumulatorpb.KeyedWindow{
		Start: timestamppb.New(time.UnixMilli(0)),
		End:   timestamppb.New(time.UnixMilli(60000)),
		Slot:  "slot-0",
		Keys:  keys,
	}

	// OPEN with a payload (no output will be produced).
	require.NoError(t, stream.Send(&accumulatorpb.AccumulatorRequest{
		Payload: &accumulatorpb.Payload{
			Keys:      keys,
			Value:     []byte("v1"),
			EventTime: timestamppb.New(time.UnixMilli(1000)),
			Watermark: timestamppb.New(time.UnixMilli(1000)),
		},
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event:       accumulatorpb.AccumulatorRequest_WindowOperation_OPEN,
			KeyedWindow: openWindow,
		},
	}), "Sending OPEN")

	closeWindow := &accumulatorpb.KeyedWindow{
		Start: timestamppb.New(time.UnixMilli(1000000)),
		End:   timestamppb.New(time.UnixMilli(2000000)),
		Slot:  "slot-7",
		Keys:  keys,
	}
	require.NoError(t, stream.Send(&accumulatorpb.AccumulatorRequest{
		Operation: &accumulatorpb.AccumulatorRequest_WindowOperation{
			Event:       accumulatorpb.AccumulatorRequest_WindowOperation_CLOSE,
			KeyedWindow: closeWindow,
		},
	}), "Sending CLOSE")

	// With no output, the very first (and only) response must be the EOF echoing the
	// CLOSE window.
	resp, err := stream.Recv()
	require.NoError(t, err, "Receiving response")
	require.True(t, resp.GetEOF(), "Expected the first response to be EOF (no output produced)")

	require.NoError(t, stream.CloseSend(), "CloseSend")

	w := resp.GetWindow()
	require.NotNil(t, w, "EOF response must carry a window")
	assert.Equal(t, closeWindow.GetEnd().GetSeconds(), w.GetEnd().GetSeconds(), "EOF window end must be the CLOSE window end, not the latest watermark")
	assert.Equal(t, closeWindow.GetStart().GetSeconds(), w.GetStart().GetSeconds(), "EOF window start")
	assert.Equal(t, "slot-7", w.GetSlot(), "EOF window slot")
	assert.Equal(t, keys, w.GetKeys(), "EOF window keys")
}
