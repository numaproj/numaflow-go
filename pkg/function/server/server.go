package server

import (
	"context"
	"fmt"
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
//	func handle(ctx context.Context, keys []string, data functionsdk.Datum) functionsdk.Messages {
//		_ = data.EventTime() // Event time is available
//		_ = data.Watermark() // Watermark is available
//		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(data.Value()))
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
//	func handle(ctx context.Context, keys []string, data functionsdk.Datum) functionsdk.MessageTs {
//		_ = data.EventTime() // Event time is available
//		_ = data.Watermark() // Watermark is available
//		return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(data.Value()).WithEventTime(time.Now()))
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
//	func handle(_ context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
//		var resultKeys = keys
//		var resultVal []byte
//		for data := range reduceCh {
//			_ = data.EventTime() // Event time is available
//			_ = data.Watermark() // Watermark is available
//		}
//		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKeys))
//	}
//
//	func main() {
//		server.New().RegisterReducer(functionsdk.ReduceFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterReducer(r functionsdk.ReduceHandler) *server {
	s.svc.Reducer = r
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr and return error.
func (s *server) Start(ctx context.Context, inputOptions ...Option) error {
	var opts = &options{
		sockAddr:       functionsdk.UDS_ADDR,
		maxMessageSize: functionsdk.DefaultMaxMessageSize,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	cleanup := func() error {
		// err if no opts.sockAddr should be ignored
		if _, err := os.Stat(opts.sockAddr); err == nil {
			return os.RemoveAll(opts.sockAddr)
		}
		return nil
	}

	if err := cleanup(); err != nil {
		return err
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	lis, err := net.Listen(functionsdk.UDS, opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", functionsdk.UDS, functionsdk.UDS_ADDR, err)
	}
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()
	functionpb.RegisterUserDefinedFunctionServer(grpcServer, s.svc)

	errCh := make(chan error, 1)
	defer close(errCh)
	// start the grpc server
	go func(ch chan<- error) {
		log.Println("starting the gRPC server with unix domain socket...", lis.Addr())
		err = grpcServer.Serve(lis)
		if err != nil {
			ch <- fmt.Errorf("failed to start the gRPC server: %v", err)
		}
	}(errCh)

	select {
	case err := <-errCh:
		return err
	case <-ctxWithSignal.Done():
		log.Println("Got a signal: terminating gRPC server...")
	}

	return nil
}
