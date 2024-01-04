package sinker

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// sinkServer is a sink gRPC server.
type sinkServer struct {
	svc  *Service
	opts *options
}

// NewServer creates a new sinkServer object.
func NewServer(h Sinker, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(sinkServer)
	s.svc = new(Service)
	s.svc.Sinker = h
	s.opts = opts
	return s
}

// Start starts the gRPC sinkServer via unix domain socket at configs.address and return error.
func (s *sinkServer) Start(ctx context.Context) error {

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

	// register the sink service
	sinkpb.RegisterSinkServer(grpcServer, s.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
