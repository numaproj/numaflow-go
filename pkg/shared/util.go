package shared

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/numaproj/numaflow-go/pkg/info"
)

const (
	uds = "unix"
)

func PrepareServer(sockAddr string, infoFilePath string) (net.Listener, error) {
	// If infoFilePath is not empty, write the server info to the file.
	if infoFilePath != "" {
		serverInfo := &info.ServerInfo{Protocol: info.UDS, Language: info.Go, MinimumNumaflowVersion: info.MinimumNumaflowVersion, Version: info.GetSDKVersion()}
		if err := info.Write(serverInfo, info.WithServerInfoFilePath(infoFilePath)); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(sockAddr); err == nil {
		if err := os.RemoveAll(sockAddr); err != nil {
			return nil, err
		}
	}

	lis, err := net.Listen(uds, sockAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to execute net.Listen(%q, %q): %v", uds, sockAddr, err)
	}

	return lis, nil
}

func CreateGRPCServer(maxMessageSize int) *grpc.Server {
	return grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
	)
}

func StartGRPCServer(ctx context.Context, grpcServer *grpc.Server, lis net.Listener) error {
	errCh := make(chan error, 1)
	defer close(errCh)
	go func(ch chan<- error) {
		log.Println("starting the gRPC server with unix domain socket...", lis.Addr())
		err := grpcServer.Serve(lis)
		if err != nil {
			ch <- fmt.Errorf("failed to start the gRPC server: %v", err)
		}
	}(errCh)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Println("Got a signal: terminating gRPC server...")
	}
	return nil
}
