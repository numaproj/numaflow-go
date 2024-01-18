package mapstreamer

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/numaproj/numaflow-go/pkg"
	mapstreampb "github.com/numaproj/numaflow-go/pkg/apis/proto/mapstream/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a map streaming gRPC server.
type server struct {
	svc  *Service
	opts *options
}

// NewServer creates a new map streaming server.
func NewServer(ms MapStreamer, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.MapperStream = ms
	s.opts = opts
	return s
}

// Start starts the map streaming gRPC server.
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
	defer grpcServer.GracefulStop()

	// register the map streaming service
	mapstreampb.RegisterMapStreamServer(grpcServer, m.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
