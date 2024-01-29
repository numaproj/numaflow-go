package sideinput

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/numaproj/numaflow-go/pkg"
	sideinputpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sideinput/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a side input gRPC server.
type server struct {
	svc  *Service
	opts *options
}

// NewSideInputServer creates a new server object.
func NewSideInputServer(r SideInputRetriever, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.Retriever = r
	s.opts = opts
	return s
}

// Start starts the gRPC server via unix domain socket at configs.address and return error.
func (s *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start listening on unix domain socket
	lis, err := shared.PrepareServer(s.opts.sockAddr, s.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := shared.CreateGRPCServer(s.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the side input service
	sideinputpb.RegisterSideInputServer(grpcServer, s.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
