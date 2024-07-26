package batchmapper

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	"github.com/numaproj/numaflow-go/pkg"
	batchmappb "github.com/numaproj/numaflow-go/pkg/apis/proto/batchmap/v1"
	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a map gRPC server.
type server struct {
	grpcServer *grpc.Server
	svc        *Service
	opts       *options
	shutdownCh <-chan struct{}
}

// NewServer creates a new batch map server.
func NewServer(m BatchMapper, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	shutdownCh := make(chan struct{})

	// create a new service and server
	svc := &Service{
		BatchMapper: m,
		shutdownCh:  shutdownCh,
	}

	return &server{
		svc:        svc,
		shutdownCh: shutdownCh,
		opts:       opts,
	}
}

// Start starts the batch map server.
func (m *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// create a server info to the file, we need to add metadata to ensure selection of the
	// correct map mode, in this case batch map
	serverInfo := info.GetDefaultServerInfo()
	serverInfo.Metadata = map[string]string{info.MapModeKey: string(info.BatchMap)}

	// start listening on unix domain socket
	lis, err := shared.PrepareServer(m.opts.sockAddr, m.opts.serverInfoFilePath, serverInfo)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	m.grpcServer = shared.CreateGRPCServer(m.opts.maxMessageSize)

	// register the batch map service
	batchmappb.RegisterBatchMapServer(m.grpcServer, m.svc)

	// start a go routine to stop the server gracefully when the context is done
	// or a shutdown signal is received from the service
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-m.shutdownCh:
		case <-ctxWithSignal.Done():
		}
		shared.StopGRPCServer(m.grpcServer)
	}()

	// start the grpc server
	if err := m.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start the gRPC server: %v", err)
	}

	// wait for the graceful shutdown to complete
	wg.Wait()
	return nil
}
