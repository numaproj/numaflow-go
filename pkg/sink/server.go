package sink

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

// sinkServer is a sink gRPC server.
type sinkServer struct {
	svc  *Service
	opts *options
}

// NewSinkServer creates a new sinkServer object.
func NewSinkServer(h SinkHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(sinkServer)
	s.svc = new(Service)
	s.svc.Sinker = h
	s.opts = opts
	return s
}

// RegisterSinker registers the sink operation handler to the sinkServer.
// Example:
//
//	func handle(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
//		result := sinksdk.ResponsesBuilder()
//		for _, datum := range datumList {
//			fmt.Println(string(datum.Value()))
//			result = result.Append(sinksdk.ResponseOK(datum.ID()))
//		}
//		return result
//	}
//
//	func main() {
//		server.NewSinkServer().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
//	}

// Start starts the gRPC sinkServer via unix domain socket at configs.SinkAddr and return error.
func (s *sinkServer) Start(ctx context.Context) error {

	// write server info to the file
	// start listening on unix domain socket
	lis, err := shared.PrepareServer(s.opts.sockAddr, s.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", shared.UDS, shared.FunctionAddr, err)
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// close the listener
	defer func() { _ = lis.Close() }()
	// create a grpc server
	grpcServer := shared.CreateGRPCServer(s.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the map service
	v1.RegisterSinkServer(grpcServer, s.svc)

	// start the grpc server
	return shared.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
