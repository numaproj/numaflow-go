package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"

	"google.golang.org/grpc"
)

type server struct {
	svc *sourcesdk.Service
}

// New creates a new server object.
func New() *server {
	s := new(server)
	s.svc = new(sourcesdk.Service)
	return s
}

// RegisterTransformer registers the map operation handler to the server.
// Example:
//
//	func transformHandle(_ context.Context, key string, d sourcesdk.Datum) sourcesdk.Messages {
//		// directly forward the input to the output without changing the event time
//		_ = d.EventTime() // Event time is available
//		_ = d.Watermark() // Watermark is available
//		return sourcesdk.MessagesBuilder().Append(sourcesdk.MessageTo(d.EventTime(), key, d.Value()))
//	}
//
//	func main() {
//		server.New().RegisterTransformer(sourcesdk.TransformFunc(transformHandle)).Start(context.Background())
//	}
func (s *server) RegisterTransformer(m sourcesdk.TransformHandler) *server {
	s.svc.Transformer = m
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr.
func (s *server) Start(ctx context.Context, inputOptions ...Option) {
	var opts = &options{
		sockAddr:       sourcesdk.Addr,
		maxMessageSize: sourcesdk.DefaultMaxMessageSize,
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

	lis, err := net.Listen(sourcesdk.Protocol, opts.sockAddr)
	if err != nil {
		log.Fatalf("failed to execute net.Listen(%q, %q): %v", sourcesdk.Protocol, sourcesdk.Addr, err)
	}
	grpcSvr := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	sourcepb.RegisterUserDefinedSourceTransformerServer(grpcSvr, s.svc)

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
