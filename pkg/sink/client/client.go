package client

import (
	"context"
	"fmt"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// client contains the grpc connection and the grpc client.
type client struct {
	conn    *grpc.ClientConn
	grpcClt sinkpb.UserDefinedSinkClient
}

// New creates a new client object.
func New(inputOptions ...Option) (*client, error) {

	var opts = &options{
		sockAddr: sinksdk.Addr,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	c := new(client)
	sockAddr := fmt.Sprintf("%s:%s", sinksdk.Protocol, opts.sockAddr)
	conn, err := grpc.Dial(sockAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to execute grpc.Dial(%q): %w", sockAddr, err)
	}
	c.conn = conn
	c.grpcClt = sinkpb.NewUserDefinedSinkClient(conn)
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

// SinkFn applies a function to a list of datum elements.
func (c *client) SinkFn(ctx context.Context, datumStreamCh <-chan *sinkpb.Datum) ([]*sinkpb.Response, error) {
	stream, err := c.grpcClt.SinkFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.SinkFn(): %w", err)
	}
	for datum := range datumStreamCh {
		if err := stream.Send(datum); err != nil {
			return nil, fmt.Errorf("failed to execute stream.Send(%v): %w", datum, err)
		}
	}
	responseList, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to execute stream.CloseAndRecv(): %w", err)
	}

	return responseList.GetResponses(), nil
}
