package server

import (
	"context"
	"log"
	"net"
	"os"

	httpfuncsdk "github.com/numaproj/numaflow-go/function"
	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
)

type server struct {
	svc *function.Service
}

// New creates a new server object.
func New() *server {
	s := new(server)
	s.svc = new(function.Service)
	return s
}

// RegisterMapper registers the map operation handler to the server.
func (s *server) RegisterMapper(m function.MapHandler) *server {
	s.svc.Mapper = m
	return s
}

// RegisterReducer registers the reduce operation handler.
func (s *server) RegisterReducer(r function.ReduceHandler) *server {
	s.svc.Reducer = r
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr.
func (s *server) Start(inputOptions ...Option) {
	var opts = &options{
		sockAddr: function.Addr,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	cleanup := func() {
		if _, err := os.Stat(opts.sockAddr); err == nil {
			if err := os.RemoveAll(opts.sockAddr); err != nil {
				log.Fatal(err)
			}
		}
	}
	cleanup()

	lis, err := net.Listen(function.Protocol, opts.sockAddr)
	if err != nil {
		log.Fatalf("failed to execute net.Listen(%q, %q): %v", function.Protocol, function.Addr, err)
	}
	grpcSvr := grpc.NewServer()
	functionpb.RegisterUserDefinedFunctionServer(grpcSvr, s.svc)
	log.Println("starting the gRPC server with unix domain socket...")
	if err := grpcSvr.Serve(lis); err != nil {
		log.Fatalf("failed to start the gRPC server: %v", err)
	}

	httpfuncsdk.Start(context.Background(), func(ctx context.Context, key, msg []byte) (httpfuncsdk.Messages, error) {
		var httpFuncSDKMessages []httpfuncsdk.Message
		msgs, err := s.svc.Mapper.HandleDo(ctx, string(key), msg)
		if err != nil {
			return httpfuncsdk.Messages{}, err
		}
		for _, m := range msgs {
			httpFuncSDKMessages = append(httpFuncSDKMessages, httpfuncsdk.Message{
				Key:   []byte(m.Key),
				Value: m.Value,
			})
		}
		return httpFuncSDKMessages, err
	})
}
