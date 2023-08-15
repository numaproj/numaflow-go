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

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/info"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
)

type server struct {
	svc *sinksdk.Service
}

// New creates a new server object.
func New() *server {
	s := new(server)
	s.svc = new(sinksdk.Service)
	return s
}

// RegisterSinker registers the sink operation handler to the server.
// Example:
//
//	func handle(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
//		result := sinksdk.ResponsesBuilder()
//		for _, datum := range datumList {
//			fmt.Println(string(datum.Value()))
//			result = result.Append(sinksdk.ResponseOK(datum.ID()))
//		}
//		return result
//	}
//
//	func main() {
//		server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterSinker(h sinksdk.SinkHandler) *server {
	s.svc.Sinker = h
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr and return error.
func (s *server) Start(ctx context.Context, inputOptions ...Option) error {
	var opts = &options{
		sockAddr:            sinksdk.Addr,
		maxMessageSize:      sinksdk.DefaultMaxMessageSize,
		sereverInfoFilePath: info.ServerInfoFilePath,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	// Write server info to the file
	serverInfo := &info.ServerInfo{Protocol: info.UDS, Language: info.Go, Version: info.GetSDKVersion()}
	if err := info.Write(serverInfo, info.WithServerInfoFilePath(opts.sereverInfoFilePath)); err != nil {
		return err
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

	lis, err := net.Listen(sinksdk.Protocol, opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", sinksdk.Protocol, sinksdk.Addr, err)
	}
	defer func() { _ = lis.Close() }()
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()
	sinkpb.RegisterUserDefinedSinkServer(grpcServer, s.svc)

	errCh := make(chan error, 1)
	defer close(errCh)
	// start the grpc server
	go func(ch chan<- error) {
		log.Println("starting the gRPC server with unix domain socket...")
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
