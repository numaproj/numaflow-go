package server

import (
	"log"
	"net"
	"os"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
)

type server struct {
	svc *function.Service
}

func NewServer() *server {
	s := new(server)
	s.svc = new(function.Service)
	return s
}

func (s *server) RegisterMapper(m function.MapHandler) *server {
	s.svc.Mapper = m
	return s
}

func (s *server) RegisterReducer(r function.ReduceHandler) *server {
	s.svc.Reducer = r
	return s
}

func (s *server) Start() {
	cleanup := func() {
		if _, err := os.Stat(function.Addr); err == nil {
			if err := os.RemoveAll(function.Addr); err != nil {
				log.Fatal(err)
			}
		}
	}
	cleanup()

	lis, err := net.Listen(function.Protocol, function.Addr)
	if err != nil {
		log.Fatalf("failed to execute net.Listen(%q, %q): %v", function.Protocol, function.Addr, err)
	}
	grpcSvr := grpc.NewServer()
	functionpb.RegisterUserDefinedFunctionServer(grpcSvr, s.svc)
	log.Println("starting the gRPC server with unix domain socket...")
	if err := grpcSvr.Serve(lis); err != nil {
		log.Fatalf("failed to start the gRPC server: %v", err)
	}
}
