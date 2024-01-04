package sourcer

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

type server struct {
	svc  *Service
	opts *options
}

// NewServer creates a new server object.
func NewServer(
	source Sourcer,
	inputOptions ...Option) numaflow.Server {
	var opts = defaultOptions()

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.Source = source
	s.opts = opts
	return s
}

// Start starts the gRPC server via unix domain socket at shared.address and return error.
func (s *server) Start(ctx context.Context) error {

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(s.opts.sockAddr, s.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// close the listener
	defer func() { _ = lis.Close() }()
	// create a grpc server
	grpcServer := shared.CreateGRPCServer(s.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the source service
	sourcepb.RegisterSourceServer(grpcServer, s.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
