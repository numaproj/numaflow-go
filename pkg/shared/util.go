package shared

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/numaproj/numaflow-go/pkg/info"
)

const (
	uds                = "unix"
	EnvUDContainerType = "NUMAFLOW_UD_CONTAINER_TYPE"
)

var ContainerType = func() string {
	if val, exists := os.LookupEnv(EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

func PrepareServer(sockAddr string, infoFilePath string, serverInfo *info.ServerInfo) (net.Listener, error) {
	// If serverInfo is not provided, then create a default server info instance.
	if serverInfo == nil {
		serverInfo = info.GetDefaultServerInfo()
	}
	// If infoFilePath is not empty, write the server info to the file.
	if infoFilePath != "" {
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

func StopGRPCServer(grpcServer *grpc.Server) {
	// Stop stops the gRPC server gracefully.
	// wait for the server to stop gracefully for 30 seconds
	// if it is not stopped, stop it forcefully
	stopped := make(chan struct{})
	go func() {
		log.Printf("gracefully stopping grpc server")
		grpcServer.GracefulStop()
		close(stopped)
	}()

	t := time.NewTimer(15 * time.Second)
	select {
	case <-t.C:
		log.Printf("forcefully stopping grpc server")
		grpcServer.Stop()
	case <-stopped:
		t.Stop()
	}
	log.Printf("grpc server stopped")
}
