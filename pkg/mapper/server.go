package mapper

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/numaproj/numaflow-go/pkg"
	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a map gRPC server.
type server struct {
	svc  *Service
	opts *options
}

// NewServer creates a new map server.
// TODO(map-batch): We call the unary map server as -> NewServer to keep backward compatibility
// Otherwise all users will have to update their current code when this change is released.
func NewServer(m Mapper, inputOptions ...Option) numaflow.Server {
	//
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.Mapper = m
	s.opts = opts
	return s
}

// NewBatchMapServer creates a new batch map server.
// TODO(map-batch): as this would be a streaming server should we see if there are some options (like maxMessageSize)
// which are different than unary server which are optimal for this use case.
func NewBatchMapServer(m BatchMapper, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.BatchMapper = m
	s.opts = opts
	return s
}

// Start starts the map server.
func (m *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(m.opts.sockAddr, m.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := shared.CreateGRPCServer(m.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the map service
	mappb.RegisterMapServer(grpcServer, m.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
