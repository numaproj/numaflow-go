package client

import (
	"context"
	"fmt"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	clientutils "github.com/numaproj/numaflow-go/pkg/sharedutils/client"
	"github.com/numaproj/numaflow-go/pkg/source"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// client contains the grpc connection and the grpc client.
type client struct {
	conn    *grpc.ClientConn
	grpcClt sourcepb.UserDefinedSourceTransformerClient
}

// New creates a new client object.
func New(inputOptions ...clientutils.Option) (*client, error) {
	var opts = &clientutils.Options{
		SockAddr: source.Addr,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	c := new(client)
	sockAddr := fmt.Sprintf("%s:%s", source.Protocol, opts.SockAddr)
	conn, err := grpc.Dial(sockAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to execute grpc.Dial(%q): %w", sockAddr, err)
	}
	c.conn = conn
	c.grpcClt = sourcepb.NewUserDefinedSourceTransformerClient(conn)
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

func (c *client) TransformFn(ctx context.Context, datum *sourcepb.Datum) ([]*sourcepb.Datum, error) {
	transformedDatumList, err := c.grpcClt.TransformFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.TransformFn(): %w", err)
	}

	return transformedDatumList.GetElements(), nil
}
