package sinker

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// sinkServer is a sink gRPC server.
type sinkServer struct {
	grpcServer *grpc.Server
	svc        *Service
	opts       *options
	shutdownCh <-chan struct{}
}

// NewServer creates a new sinkServer object.
func NewServer(h Sinker, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	shutdownCh := make(chan struct{})

	// create a new service and server
	svc := &Service{
		Sinker:     h,
		shutdownCh: shutdownCh,
	}

	return &sinkServer{
		svc:        svc,
		shutdownCh: shutdownCh,
		opts:       opts,
	}
}

// Start starts the gRPC sinkServer via unix domain socket at configs.address and return error.
func (s *sinkServer) Start(ctx context.Context) error {

	// write server info to the file
	serverInfo := info.GetDefaultServerInfo()
	serverInfo.MinimumNumaflowVersion = info.MinimumNumaflowVersion[info.Sinker]
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(s.opts.sockAddr, s.opts.serverInfoFilePath, serverInfo)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	s.grpcServer = shared.CreateGRPCServer(s.opts.maxMessageSize)

	// register the sink service
	sinkpb.RegisterSinkServer(s.grpcServer, s.svc)

	// start a go routine to stop the server gracefully when the context is done
	// or a shutdown signal is received from the service
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-s.shutdownCh:
			log.Printf("shutdown signal received")
		case <-ctxWithSignal.Done():
		}
		shared.StopGRPCServer(s.grpcServer)
	}()

	// start the grpc server
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start the gRPC server: %v", err)
	}

	// wait for the graceful shutdown to complete
	wg.Wait()
	return nil
}
