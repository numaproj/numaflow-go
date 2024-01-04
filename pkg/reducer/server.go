package reducer

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	reducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a reduce gRPC server.
type server struct {
	svc  *Service
	opts *options
}

// NewServer creates a new reduce server.
func NewServer(r ReducerCreator, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.reducerCreatorHandle = r
	s.opts = opts
	return s
}

// Start starts the reduce gRPC server.
func (r *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(r.opts.sockAddr, r.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := shared.CreateGRPCServer(r.opts.maxMessageSize)
	defer grpcServer.GracefulStop()

	// register the reduce service
	reducepb.RegisterReduceServer(grpcServer, r.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
