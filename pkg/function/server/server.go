package function

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
)

func NewServer() *server {
	s := new(server)
	s.svc = new(service)
	return s
}

const (
	protocol = "unix"
	sockAddr = "/tmp/uds/numaflow.sock"
)

type server struct {
	svc *service
}

func (s *server) RegisterMapper(m MapHandler) *server {
	s.svc.mapper = m
	return s
}

func (s *server) RegisterReducer(r ReduceHandler) *server {
	s.svc.reducer = r
	return s
}

func (s *server) Start() {
	cleanup := func() {
		if _, err := os.Stat(sockAddr); err == nil {
			if err := os.RemoveAll(sockAddr); err != nil {
				log.Fatal(err)
			}
		}
	}

	cleanup()

	lis, err := net.Listen(protocol, sockAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcSvr := grpc.NewServer()
	functionpb.RegisterUserDefinedFunctionServer(grpcSvr, s.svc)
	log.Println("Starting gRPC server with unix domain socket...")

	if err := grpcSvr.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
