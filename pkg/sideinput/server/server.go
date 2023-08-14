package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	sideinputpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sideinput/v1"
	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

type server struct {
	svc *sideinputsdk.Service
}

// NewSideInputServer creates a new server object.
func NewSideInputServer(handler sideinputsdk.RetrieveSideInputHandler) *server {
	s := new(server)
	s.svc = new(sideinputsdk.Service)
	s.svc.Retriever = handler
	return s
}

// RegisterRetriever registers the retrieve side input operation handler to the server.
// Example:
//
//	func handle(ctx context.Context) sideinputsdk.MessageSI {
//		return sideinputsdk.NewMessageSI(sideInput.Value()))
//	}
//
//	func main() {
//		server.NewSideInputServer().RegisterRetriever(sideinputsdk.RetrieveSideInputFunc(handle)).Start(context.Background())
//	}
func (s *server) RegisterRetriever(m sideinputsdk.RetrieveSideInputHandler) *server {
	s.svc.Retriever = m
	return s
}

// Start starts the gRPC server via unix domain socket at configs.Addr and return error.
func (s *server) Start(ctx context.Context, inputOptions ...Option) error {
	var opts = &options{
		sockAddr:       sideinputsdk.Addr,
		maxMessageSize: sideinputsdk.DefaultMaxMessageSize,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	cleanup := func() error {
		// err if no opts.sockAddr should be ignored
		if _, err := os.Stat(opts.sockAddr); err == nil {
			return os.RemoveAll(opts.sockAddr)
		}
		return nil
	}

	if err := cleanup(); err != nil {
		return err
	}

	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	lis, err := net.Listen(sideinputsdk.Protocol, opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", sideinputsdk.Protocol, sideinputsdk.Addr, err)
	}
	defer func() { _ = lis.Close() }()
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(opts.maxMessageSize),
		grpc.MaxSendMsgSize(opts.maxMessageSize),
	)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()
	sideinputpb.RegisterUserDefinedSideInputServer(grpcServer, s.svc)

	errCh := make(chan error, 1)
	defer close(errCh)
	// start the grpc server
	go func(ch chan<- error) {
		log.Println("starting the gRPC server with unix domain socket...", lis.Addr())
		err = grpcServer.Serve(lis)
		if err != nil {
			ch <- fmt.Errorf("failed to start the gRPC server: %v", err)
		}
	}(errCh)

	select {
	case err := <-errCh:
		return err
	case <-ctxWithSignal.Done():
		log.Println("Got a signal: terminating gRPC server...")
	}

	return nil
}
