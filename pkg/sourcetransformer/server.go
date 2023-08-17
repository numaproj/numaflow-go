package sourcetransformer

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/numaproj/numaflow-go/pkg"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

type server struct {
	svc  *Service
	opts *options
}

// NewServer creates a new SourceTransformer server.
func NewServer(m SourceTransformer, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(server)
	s.svc = new(Service)
	s.svc.Transformer = m
	s.opts = opts
	return s
}

// Start starts the SourceTransformer server.
func (m *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(m.opts.sockAddr, m.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", shared.UDS, shared.SourceTransformerAddr, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := shared.CreateGRPCServer(m.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the map service
	v1.RegisterSourceTransformServer(grpcServer, m.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
