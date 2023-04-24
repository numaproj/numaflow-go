package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/emptypb"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/info"
)

// client contains the grpc connection and the grpc client.
type client struct {
	conn    *grpc.ClientConn
	grpcClt functionpb.UserDefinedFunctionClient
}

// New creates a new client object.
func New(inputOptions ...Option) (*client, error) {
	var opts = &options{
		maxMessageSize:             function.DefaultMaxMessageSize,
		serverInfoFilePath:         info.ServerInfoFilePath,
		tcpSockAddr:                function.TCP_ADDR,
		udsSockAddr:                function.UDS_ADDR,
		serverInfoReadinessTimeout: 120 * time.Second, // Default timeout is 120 seconds
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.serverInfoReadinessTimeout)
	defer cancel()

	if err := info.WaitUntilReady(ctx, info.WithServerInfoFilePath(opts.serverInfoFilePath)); err != nil {
		return nil, fmt.Errorf("failed to wait until server info is ready: %w", err)
	}

	serverInfo, err := info.Read(info.WithServerInfoFilePath(opts.serverInfoFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to read server info: %w", err)
	}
	// TODO: Use serverInfo to check compatibility.
	if serverInfo != nil {
		log.Printf("ServerInfo: %v\n", serverInfo)
	}

	c := new(client)
	var conn *grpc.ClientConn
	var sockAddr string
	// Make a TCP connection client for multiprocessing grpc server
	if serverInfo.Protocol == info.TCP {
		// Populate connection variables for client connection
		// based on multiprocessing enabled/disabled
		if err := regMultProcResolver(serverInfo); err != nil {
			return nil, fmt.Errorf("failed to start Multiproc Client: %w", err)
		}

		sockAddr = fmt.Sprintf("%s%s", connAddr, opts.tcpSockAddr)
		log.Println("Multiprocessing TCP Client:", sockAddr)
		conn, err = grpc.Dial(
			fmt.Sprintf("%s:///%s", custScheme, custServiceName),
			// This sets the initial load balancing policy as Round Robin
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(opts.maxMessageSize), grpc.MaxCallSendMsgSize(opts.maxMessageSize)),
		)
	} else {
		sockAddr = fmt.Sprintf("%s:%s", function.UDS, opts.udsSockAddr)
		log.Println("UDS Client:", sockAddr)
		conn, err = grpc.Dial(sockAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(opts.maxMessageSize), grpc.MaxCallSendMsgSize(opts.maxMessageSize)))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute grpc.Dial(%q): %w", sockAddr, err)
	}
	c.conn = conn
	c.grpcClt = functionpb.NewUserDefinedFunctionClient(conn)
	return c, nil
}

// CloseConn closes the grpc client connection.
func (c *client) CloseConn(ctx context.Context) error {
	return c.conn.Close()
}

// IsReady returns true if the grpc connection is ready to use.
func (c *client) IsReady(ctx context.Context, in *emptypb.Empty) (bool, error) {
	resp, err := c.grpcClt.IsReady(ctx, in)
	if err != nil {
		return false, err
	}
	return resp.GetReady(), nil
}

// MapFn applies a function to each datum element.
func (c *client) MapFn(ctx context.Context, datum *functionpb.DatumRequest) ([]*functionpb.DatumResponse, error) {
	mappedDatumList, err := c.grpcClt.MapFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.MapFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
// MapTFn can be used only at source vertex by source data transformer.
func (c *client) MapTFn(ctx context.Context, datum *functionpb.DatumRequest) ([]*functionpb.DatumResponse, error) {
	mappedDatumList, err := c.grpcClt.MapTFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.MapTFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// ReduceFn applies a reduce function to a datum stream.
func (c *client) ReduceFn(ctx context.Context, datumStreamCh <-chan *functionpb.DatumRequest) ([]*functionpb.DatumResponse, error) {
	var g errgroup.Group
	datumList := make([]*functionpb.DatumResponse, 0)

	stream, err := c.grpcClt.ReduceFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.ReduceFn(): %w", err)
	}
	// stream the messages to server
	g.Go(func() error {
		var sendErr error
		for datum := range datumStreamCh {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if sendErr = stream.Send(datum); sendErr != nil {
					// we don't need to invoke close on the stream
					// if there is an error gRPC will close the stream.
					return sendErr
				}
			}
		}
		return stream.CloseSend()
	})

	// read the response from the server stream
outputLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var resp *functionpb.DatumResponseList
			resp, err = stream.Recv()
			if err == io.EOF {
				break outputLoop
			}
			if err != nil {
				return nil, err
			}
			datumList = append(datumList, resp.Elements...)
		}
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return datumList, err
}

// setConn function is used to populate the connection properties based
// on multiprocessing TCP or UDS connection

func regMultProcResolver(svrInfo *info.ServerInfo) error {
	numCpu, err := strconv.Atoi(svrInfo.Metadata["CPU_LIMIT"])
	if err != nil {
		return err
	}
	log.Println("Num CPU:", numCpu)
	conn := buildConnAddrs(numCpu)
	res := &multiProcResolverBuilder{addrsList: conn}
	resolver.Register(res)
	log.Println("TCP client list:", res.addrsList)
	return nil
}
