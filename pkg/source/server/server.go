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

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/info"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
)

type server struct {
	svc *sourcesdk.Service
}

// New creates a new server object.
func New(
	pendingHandler sourcesdk.PendingHandler,
	readHandler sourcesdk.ReadHandler,
) *server {
	s := new(server)
	s.svc = new(sourcesdk.Service)
	s.svc.PendingHandler = pendingHandler
	s.svc.ReadHandler = readHandler
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr and return error.
func (s *server) Start(ctx context.Context, inputOptions ...Option) error {
	var opts = &options{
		sockAddr:           sourcesdk.UdsAddr,
		maxMessageSize:     sourcesdk.DefaultMaxMessageSize,
		serverInfoFilePath: info.ServerInfoFilePath,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	// Write server info to the file
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
	lis, err := net.Listen(sourcesdk.UDS, opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", sourcesdk.UDS, sourcesdk.UdsAddr, err)
	}
	defer func() { _ = lis.Close() }()
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()
	sourcepb.RegisterSourceServer(grpcServer, s.svc)

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
