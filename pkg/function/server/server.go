package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
)

type server struct {
	svc *functionsdk.Service
}

// New creates a new server object.
func New() *server {
	s := new(server)
	s.svc = new(functionsdk.Service)
	return s
}

// RegisterMapper registers the map operation handler to the server.
// Example:
//
//	func handle(ctx context.Context, key string, data functionsdk.Datum) functionsdk.Messages {
//		_ = data.EventTime() // Event time is available
//		_ = data.Watermark() // Watermark is available
//		return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(data.Value()))
//	}
//
//	func main() {
//		server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterMapper(m functionsdk.MapHandler) *server {
	s.svc.Mapper = m
	return s
}

// RegisterMapperT registers the mapT operation handler to the server.
// Example:
//
//	func handle(ctx context.Context, key string, data functionsdk.Datum) functionsdk.MessageTs {
//		_ = data.EventTime() // Event time is available
//		_ = data.Watermark() // Watermark is available
//		return functionsdk.MessageTsBuilder().Append(functionsdk.MessageTToAll(time.Now(), data.Value()))
//	}
//
//	func main() {
//		server.New().RegisterMapperT(functionsdk.MapTFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterMapperT(m functionsdk.MapTHandler) *server {
	s.svc.MapperT = m
	return s
}

// RegisterReducer registers the reduce operation handler.
// Example:
//
//	func handle(_ context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
//		var resultKey = key
//		var resultVal []byte
//		for data := range reduceCh {
//			_ = data.EventTime() // Event time is available
//			_ = data.Watermark() // Watermark is available
//		}
//		return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal))
//	}
//
//	func main() {
//		server.New().RegisterReducer(functionsdk.ReduceFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterReducer(r functionsdk.ReduceHandler) *server {
	s.svc.Reducer = r
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr.
func (s *server) Start(ctx context.Context, inputOptions ...Option) {
	var opts = &options{
		sockAddr:       functionsdk.Addr,
		maxMessageSize: functionsdk.DefaultMaxMessageSize,
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

	lis, err := net.Listen(functionsdk.Protocol, opts.sockAddr)
	if err != nil {
		log.Fatalf("failed to execute net.Listen(%q, %q): %v", functionsdk.Protocol, functionsdk.Addr, err)
	}
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	functionpb.RegisterUserDefinedFunctionServer(grpcServer, s.svc)

	// start the grpc server
	go func() {
		log.Println("starting the gRPC server with unix domain socket...")
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to start the gRPC server: %v", err)
		}
	}()

	<-ctxWithSignal.Done()
	log.Println("Got a signal: terminating gRPC server...")
	defer log.Println("Successfully stopped the gRPC server")
	grpcServer.GracefulStop()
}
