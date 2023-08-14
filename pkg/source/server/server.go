package server

import (
	"context"
	"fmt"
	"log"

	"github.com/numaproj/numaflow-go/pkg"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/source/transformerfn"
	"github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/mapt"
	"github.com/numaproj/numaflow-go/pkg/util"
)

type SourceTransformerServer struct {
	svc  *mapt.Service
	opts *options
}

// NewSourceTransformerServer creates a new map server.
func NewSourceTransformerServer(ctx context.Context, m source.MapTHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(SourceTransformerServer)
	s.svc = new(mapt.Service)
	s.svc.MapperT = m
	s.opts = opts
	return s
}

func (m *SourceTransformerServer) Start(ctx context.Context) error {
	// write server info to the file
	// start listening on unix domain socket
	lis, err := util.PrepareServer(m.opts.sockAddr, m.opts.serverInfoFilePath)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", util.UDS, util.FunctionAddr, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := util.CreateGRPCServer(m.opts.maxMessageSize)
	defer log.Println("Successfully stopped the gRPC server")
	defer grpcServer.GracefulStop()

	// register the map service
	transformerfn.RegisterSourceTransformerServer(grpcServer, m.svc)

	// start the grpc server
	return util.StartGRPCServer(ctx, grpcServer, lis)
}
