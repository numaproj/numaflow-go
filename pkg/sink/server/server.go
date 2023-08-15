package server

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	numaflow "github.com/numaproj/numaflow-go/pkg"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/sinkfn"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/util"
)

// sinkServer is a sink gRPC server.
type sinkServer struct {
	svc  *sinksdk.Service
	opts *options
}

// NewSinkServer creates a new sinkServer object.
func NewSinkServer(h sinksdk.SinkHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(sinkServer)
	s.svc = new(sinksdk.Service)
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
	lis, err := util.PrepareServer(s.opts.sockAddr, s.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", util.UDS, util.FunctionAddr, err)
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// close the listener
	defer func() { _ = lis.Close() }()
	// create a grpc server
	grpcServer := util.CreateGRPCServer(s.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the map service
	sinkfn.RegisterSinkServer(grpcServer, s.svc)

	// start the grpc server
	return util.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
