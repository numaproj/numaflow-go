package sourcetransformer

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"

	proto "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newServer(t *testing.T, register func(server *grpc.Server)) *grpc.ClientConn {
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

var testTime = time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)

func TestService_sourceTransformFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *proto.SourceTransformRequest
	}

	tests := []struct {
		name    string
		handler SourceTransformer
		args    args
		want    *proto.SourceTransformResponse
	}{
		{
			name: "sourceTransform_fn_forward_msg",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg, testTime).WithKeys([]string{keys[0] + "_test"}))
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.SourceTransformRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			want: &proto.SourceTransformResponse{
				Results: []*proto.SourceTransformResponse_Result{
					{
						EventTime: timestamppb.New(testTime),
						Keys:      []string{"client_test"},
						Value:     []byte(`test`),
					},
				},
			},
		},
		{
			name: "sourceTransform_fn_forward_msg_forward_to_all",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg, testTime))
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.SourceTransformRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			want: &proto.SourceTransformResponse{
				Results: []*proto.SourceTransformResponse_Result{
					{
						EventTime: timestamppb.New(testTime),
						Value:     []byte(`test`),
					},
				},
			},
		},
		{
			name: "sourceTransform_fn_forward_msg_drop_msg",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				return MessagesBuilder().Append(MessageToDrop(testTime))
			}),
			args: args{
				ctx: context.Background(),
				d: &proto.SourceTransformRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: timestamppb.New(time.Time{}),
					Watermark: timestamppb.New(time.Time{}),
				},
			},
			want: &proto.SourceTransformResponse{
				Results: []*proto.SourceTransformResponse_Result{
					{
						EventTime: timestamppb.New(testTime),
						Tags:      []string{DROP},
						Value:     nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				Transformer: tt.handler,
			}

			conn := newServer(t, func(server *grpc.Server) {
				proto.RegisterSourceTransformServer(server, svc)
			})

			client := proto.NewSourceTransformClient(conn)
			stream, err := client.SourceTransformFn(context.Background())
			assert.NoError(t, err, "Creating stream")

			err = stream.Send(tt.args.d)
			assert.NoError(t, err, "Sending message over the stream")

			got, err := stream.Recv()
			assert.NoError(t, err, "Receiving message from the stream")

			assert.Equal(t, got.Results, tt.want.Results)
		})
	}
}

func TestService_SourceTransformFn_Multiple_Messages(t *testing.T) {
	svc := &Service{
		Transformer: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
			msg := datum.Value()
			return MessagesBuilder().Append(NewMessage(msg, testTime).WithKeys([]string{keys[0] + "_test"}))
		}),
	}
	conn := newServer(t, func(server *grpc.Server) {
		proto.RegisterSourceTransformServer(server, svc)
	})

	client := proto.NewSourceTransformClient(conn)
	stream, err := client.SourceTransformFn(context.Background())
	assert.NoError(t, err, "Creating stream")

	const msgCount = 10
	for i := 0; i < msgCount; i++ {
		msg := proto.SourceTransformRequest{
			Keys:      []string{"client"},
			Value:     []byte(fmt.Sprintf("test_%d", i)),
			EventTime: timestamppb.New(time.Time{}),
			Watermark: timestamppb.New(time.Time{}),
		}
		err = stream.Send(&msg)
		assert.NoError(t, err, "Sending message over the stream")
	}
	err = stream.CloseSend()
	assert.NoError(t, err, "Closing the send direction of the stream")

	expectedResults := make([][]*proto.SourceTransformResponse_Result, msgCount)
	for i := 0; i < msgCount; i++ {
		expectedResults[i] = []*proto.SourceTransformResponse_Result{
			{
				EventTime: timestamppb.New(testTime),
				Keys:      []string{"client_test"},
				Value:     []byte(fmt.Sprintf("test_%d", i)),
			},
		}
	}

	results := make([][]*proto.SourceTransformResponse_Result, msgCount)
	for i := 0; i < msgCount; i++ {
		got, err := stream.Recv()
		assert.NoError(t, err, "Receiving message from the stream")
		results[i] = got.Results
	}
	assert.ElementsMatch(t, results, expectedResults)
}
