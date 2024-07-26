package reducer

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	reducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// server is a reduce gRPC server.
type server struct {
	grpcServer *grpc.Server
	svc        *Service
	opts       *options
	shutdownCh <-chan struct{}
}

// NewServer creates a new reduce server.
func NewServer(r ReducerCreator, inputOptions ...Option) numaflow.Server {
	opts := defaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	shutdownCh := make(chan struct{})

	// create a new service and server
	svc := &Service{
		reducerCreatorHandle: r,
		shutdownCh:           shutdownCh,
	}

	return &server{
		svc:        svc,
		shutdownCh: shutdownCh,
		opts:       opts,
	}
}

// Start starts the reduce gRPC server.
func (r *server) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(r.opts.sockAddr, r.opts.serverInfoFilePath, info.GetDefaultServerInfo())
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, address, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	r.grpcServer = shared.CreateGRPCServer(r.opts.maxMessageSize)
	defer r.grpcServer.GracefulStop()

	// register the reduce service
	reducepb.RegisterReduceServer(r.grpcServer, r.svc)

	// start a go routine to stop the server gracefully
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-r.shutdownCh:
		case <-ctxWithSignal.Done():
		}
		shared.StopGRPCServer(r.grpcServer)
	}()

	// start the grpc server
	if err := r.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start the gRPC server: %v", err)
	}

	// wait for the graceful shutdown to complete
	wg.Wait()
	return nil
}
