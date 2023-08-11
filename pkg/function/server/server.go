package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/info"
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
// See an example at pkg/function/examples/flatmap/main.go
func (s *server) RegisterMapper(m functionsdk.MapHandler) *server {
	s.svc.Mapper = m
	return s
}

// RegisterMapperStream registers the mapStream operation handler to the server.
// See an example at pkg/function/examples/flatmap_stream/main.go
func (s *server) RegisterMapperStream(m functionsdk.MapStreamHandler) *server {
	s.svc.MapperStream = m
	return s
}

// RegisterMapperT registers the mapT operation handler to the server.
// See an example at pkg/function/examples/assign_event_time/main.go
func (s *server) RegisterMapperT(m functionsdk.MapTHandler) *server {
	s.svc.MapperT = m
	return s
}

// RegisterReducer registers the reduce operation handler.
// See an example at pkg/function/examples/sum/main.go
func (s *server) RegisterReducer(r functionsdk.ReduceHandler) *server {
	s.svc.Reducer = r
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr and returns error if any.
func (s *server) Start(ctx context.Context, inputOptions ...Option) error {
	var opts = &options{
		sockAddr:           functionsdk.UdsAddr,
		maxMessageSize:     functionsdk.DefaultMaxMessageSize,
		serverInfoFilePath: info.ServerInfoFilePath,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	// Write server info to the file.
	serverInfo := &info.ServerInfo{Protocol: info.UDS, Language: info.Go, Version: info.GetSDKVersion()}
	if err := info.Write(serverInfo, info.WithServerInfoFilePath(opts.serverInfoFilePath)); err != nil {
		return err
	}

	// cleanup cleans up the unix domain socket file if exists.
	// TODO - once we support TCP, we need to update it to support TCP as well.
	cleanup := func() error {
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
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", functionsdk.UDS, functionsdk.UdsAddr, err)
	}
	defer func() { _ = lis.Close() }()
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
