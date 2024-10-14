package mapstreamer

import (
	"context"
	"errors"
	"fmt"
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

	proto "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestService_MapFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *proto.MapRequest
	}

	tests := []struct {
		name    string
		handler MapStreamer
		args    args
		want    *proto.MapResponse
	}{
		{
			name: "map_fn_forward_msg",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				defer close(messageCh)
				msg := datum.Value()
				messageCh <- NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.MapRequest{
					Request: &proto.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
				},
			},
			want: &proto.MapResponse{
				Results: []*proto.MapResponse_Result{
					{
						Keys:  []string{"client_test"},
						Value: []byte(`test`),
					},
				},
			},
		},
		{
			name: "map_fn_forward_msg_forward_to_all",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				defer close(messageCh)
				msg := datum.Value()
				messageCh <- NewMessage(msg)
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.MapRequest{
					Request: &proto.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
				},
			},
			want: &proto.MapResponse{
				Results: []*proto.MapResponse_Result{
					{
						Value: []byte(`test`),
					},
				},
			},
		},
		{
			name: "map_fn_forward_msg_drop_msg",
			handler: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
				defer close(messageCh)
				messageCh <- MessageToDrop()
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.MapRequest{
					Request: &proto.MapRequest_Request{
						Keys:      []string{"client"},
						Value:     []byte(`test`),
						EventTime: timestamppb.New(time.Time{}),
						Watermark: timestamppb.New(time.Time{}),
					},
				},
			},
			want: &proto.MapResponse{
				Results: []*proto.MapResponse_Result{
					{
						Tags:  []string{DROP},
						Value: nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				MapperStream: tt.handler,
			}

			conn := newTestServer(t, func(server *grpc.Server) {
				proto.RegisterMapServer(server, svc)
			})

			client := proto.NewMapClient(conn)
			stream, err := client.MapFn(context.Background())
			require.NoError(t, err, "Creating stream")

			doHandshake(t, stream)

			err = stream.Send(tt.args.d)
			require.NoError(t, err, "Sending message over the stream")

			got, err := stream.Recv()
			require.NoError(t, err, "Receiving message from the stream")

			assert.Equal(t, got.Results, tt.want.Results)
		})
	}
}

func doHandshake(t *testing.T, stream proto.Map_MapFnClient) {
	t.Helper()
	handshakeReq := &proto.MapRequest{
		Handshake: &proto.Handshake{Sot: true},
	}
	err := stream.Send(handshakeReq)
	require.NoError(t, err, "Sending handshake request to the stream")

	handshakeResp, err := stream.Recv()
	require.NoError(t, err, "Receiving handshake response")

	require.Empty(t, handshakeResp.Results, "Invalid handshake response")
	require.Empty(t, handshakeResp.Id, "Invalid handshake response")
	require.NotNil(t, handshakeResp.Handshake, "Invalid handshake response")
	require.True(t, handshakeResp.Handshake.Sot, "Invalid handshake response")
}

func TestService_MapFn_SingleMessage_MultipleResponses(t *testing.T) {
	svc := &Service{
		MapperStream: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
			defer close(messageCh)
			for i := 0; i < 10; i++ {
				msg := fmt.Sprintf("response_%d", i)
				messageCh <- NewMessage([]byte(msg)).WithKeys([]string{keys[0] + "_test"})
			}
		}),
	}
	conn := newTestServer(t, func(server *grpc.Server) {
		proto.RegisterMapServer(server, svc)
	})

	client := proto.NewMapClient(conn)
	stream, err := client.MapFn(context.Background())
	require.NoError(t, err, "Creating stream")

	doHandshake(t, stream)

	msg := proto.MapRequest{
		Request: &proto.MapRequest_Request{
			Keys:      []string{"client"},
			Value:     []byte("test"),
			EventTime: timestamppb.New(time.Time{}),
			Watermark: timestamppb.New(time.Time{}),
		},
	}
	err = stream.Send(&msg)
	require.NoError(t, err, "Sending message over the stream")
	err = stream.CloseSend()
	require.NoError(t, err, "Closing the send direction of the stream")

	expectedResults := make([][]*proto.MapResponse_Result, 10)
	for i := 0; i < 10; i++ {
		expectedResults[i] = []*proto.MapResponse_Result{
			{
				Keys:  []string{"client_test"},
				Value: []byte(fmt.Sprintf("response_%d", i)),
			},
		}
	}

	results := make([][]*proto.MapResponse_Result, 0)
	for i := 0; i < 10; i++ {
		got, err := stream.Recv()
		require.NoError(t, err, "Receiving message from the stream")
		results = append(results, got.Results)
	}

	require.ElementsMatch(t, results, expectedResults)

}

func TestService_MapFn_Multiple_Messages(t *testing.T) {
	svc := &Service{
		MapperStream: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
			defer close(messageCh)
			msg := datum.Value()
			messageCh <- NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
		}),
	}
	conn := newTestServer(t, func(server *grpc.Server) {
		proto.RegisterMapServer(server, svc)
	})

	client := proto.NewMapClient(conn)
	stream, err := client.MapFn(context.Background())
	require.NoError(t, err, "Creating stream")

	doHandshake(t, stream)

	const msgCount = 10
	for i := 0; i < msgCount; i++ {
		msg := proto.MapRequest{
			Request: &proto.MapRequest_Request{
				Keys:      []string{"client"},
				Value:     []byte(fmt.Sprintf("test_%d", i)),
				EventTime: timestamppb.New(time.Time{}),
				Watermark: timestamppb.New(time.Time{}),
			},
		}
		err = stream.Send(&msg)
		require.NoError(t, err, "Sending message over the stream")
	}
	err = stream.CloseSend()
	require.NoError(t, err, "Closing the send direction of the stream")

	expectedResults := make([][]*proto.MapResponse_Result, msgCount)
	results := make([][]*proto.MapResponse_Result, 0)

	for i := 0; i < msgCount; i++ {
		expectedResults[i] = []*proto.MapResponse_Result{
			{
				Keys:  []string{"client_test"},
				Value: []byte(fmt.Sprintf("test_%d", i)),
			},
		}

		got, err := stream.Recv()
		require.NoError(t, err, "Receiving message from the stream")
		results = append(results, got.Results)

		eot, err := stream.Recv()
		require.NoError(t, err, "Receiving message from the stream")
		require.True(t, eot.Status.Eot, "Expected EOT message")
	}

	require.ElementsMatch(t, results, expectedResults)
}

func TestService_MapFn_Panic(t *testing.T) {
	svc := &Service{
		MapperStream: MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
			panic("map failed")
		}),
		shutdownCh: make(chan<- struct{}, 1),
	}
	conn := newTestServer(t, func(server *grpc.Server) {
		proto.RegisterMapServer(server, svc)
	})

	client := proto.NewMapClient(conn)
	stream, err := client.MapFn(context.Background())
	require.NoError(t, err, "Creating stream")

	doHandshake(t, stream)

	msg := proto.MapRequest{
		Request: &proto.MapRequest_Request{
			Keys:      []string{"client"},
			Value:     []byte("test"),
			EventTime: timestamppb.New(time.Time{}),
			Watermark: timestamppb.New(time.Time{}),
		},
	}
	err = stream.Send(&msg)
	require.NoError(t, err, "Sending message over the stream")
	err = stream.CloseSend()
	require.NoError(t, err, "Closing the send direction of the stream")
	_, err = stream.Recv()
	require.Error(t, err, "Expected error while receiving message from the stream")
	gotStatus, _ := status.FromError(err)
	expectedStatus := status.Convert(status.Errorf(codes.Internal, "error processing requests: panic inside mapStream handler: map failed"))
	require.Equal(t, expectedStatus, gotStatus)
}
