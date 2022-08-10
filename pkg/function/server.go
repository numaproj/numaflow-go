package function

import (
	"log"
	"net"

	"google.golang.org/grpc"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
)

func Initialize() *srv {
	s := new(srv)
	s.svc = new(functionService)
	return s
}

const (
	port = ":7777"
)

type srv struct {
	svc *functionService
}

func (s *srv) RegisterMapper(m MapHandler) *srv {
	s.svc.mapper = m
	return s
}

func (s *srv) RegisterReducer(r ReduceHandler) *srv {
	s.svc.reducer = r
	return s
}

func (s *srv) Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcSvr := grpc.NewServer()
	v1.RegisterFunctionServiceServer(grpcSvr, s.svc)
	log.Printf("Starting gRPC server at %s\n", port)
	if err := grpcSvr.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
