package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/sink"
	"google.golang.org/grpc"
)

type server struct {
	svc *sink.Service
}

// New creates a new server object.
func New() *server {
	s := new(server)
	s.svc = new(sink.Service)
	return s
}

// RegisterSinker registers the sink operation handler to the server.
func (s *server) RegisterSinker(h sink.SinkHandler) *server {
	s.svc.Sinker = h
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr.
func (s *server) Start(ctx context.Context, inputOptions ...Option) {
	var opts = &options{
		sockAddr:       sink.Addr,
		maxMessageSize: sink.DefaultMaxMessageSize,
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

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen(sink.Protocol, opts.sockAddr)
	if err != nil {
		log.Fatalf("failed to execute net.Listen(%q, %q): %v", sink.Protocol, sink.Addr, err)
	}
	grpcSvr := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	sinkpb.RegisterUserDefinedSinkServer(grpcSvr, s.svc)

	// start the grpc server
	go func() {
		log.Println("starting the gRPC server with unix domain socket...")
		err = grpcSvr.Serve(lis)
		if err != nil {
			log.Fatalf("failed to start the gRPC server: %v", err)
		}
	}()

	<-ctxWithSignal.Done()
	log.Println("Got a signal: terminating gRPC server...")
	defer log.Println("Successfully stopped the gRPC server")
	grpcSvr.GracefulStop()
}
