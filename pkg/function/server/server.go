package server

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/numaproj/numaflow-go/pkg"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/mapfn"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/reducefn"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/smapfn"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	mapsvc "github.com/numaproj/numaflow-go/pkg/function/map"
	"github.com/numaproj/numaflow-go/pkg/function/reduce"
	"github.com/numaproj/numaflow-go/pkg/function/smap"
	"github.com/numaproj/numaflow-go/pkg/util"
)

type mapServer struct {
	svc  *mapsvc.Service
	opts *options
}

// NewMapServer creates a new map server.
func NewMapServer(ctx context.Context, m functionsdk.MapHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(mapServer)
	s.svc = new(mapsvc.Service)
	s.svc.Mapper = m
	s.opts = opts
	return s
}

func (m *mapServer) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
	mapfn.RegisterMapServer(grpcServer, m.svc)

	// start the grpc server
	return util.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}

type reduceServer struct {
	svc  *reduce.Service
	opts *options
}

// NewReduceServer creates a new reduce server.
func NewReduceServer(ctx context.Context, r functionsdk.ReduceHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(reduceServer)
	s.svc = new(reduce.Service)
	s.svc.Reducer = r
	s.opts = opts
	return s
}

func (r *reduceServer) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := util.PrepareServer(r.opts.serverInfoFilePath, r.opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", util.UDS, util.FunctionAddr, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := util.CreateGRPCServer(r.opts.maxMessageSize)
	defer grpcServer.GracefulStop()

	// register the reduce service
	reducefn.RegisterReduceServer(grpcServer, r.svc)

	// start the grpc server
	return util.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}

type mapStreamServer struct {
	svc  *smap.Service
	opts *options
}

// NewMapStreamServer creates a new map streaming server.
func NewMapStreamServer(ctx context.Context, ms functionsdk.MapStreamHandler, inputOptions ...Option) numaflow.Server {
	opts := DefaultOptions()
	for _, inputOption := range inputOptions {
		inputOption(opts)
	}
	s := new(mapStreamServer)
	s.svc = new(smap.Service)
	s.svc.MapperStream = ms
	s.opts = opts
	return s
}

func (m *mapStreamServer) Start(ctx context.Context) error {
	ctxWithSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// write server info to the file
	// start listening on unix domain socket
	lis, err := util.PrepareServer(m.opts.serverInfoFilePath, m.opts.sockAddr)
	if err != nil {
		return fmt.Errorf("failed to execute net.Listen(%q, %q): %v", util.UDS, util.FunctionAddr, err)
	}
	// close the listener
	defer func() { _ = lis.Close() }()

	// create a grpc server
	grpcServer := util.CreateGRPCServer(m.opts.maxMessageSize)
	defer grpcServer.GracefulStop()

	// register the map streaming service
	smapfn.RegisterMapStreamServer(grpcServer, m.svc)

	// start the grpc server
	return util.StartGRPCServer(ctxWithSignal, grpcServer, lis)
}
