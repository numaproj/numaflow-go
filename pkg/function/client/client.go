package client

import (
	"context"
	"fmt"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Client contains the grpc connection and the grpc client.
type Client struct {
	conn    *grpc.ClientConn
	grpcClt functionpb.UserDefinedFunctionClient
}

// New creates a new client object.
func New(inputOptions ...Option) (*Client, error) {
	var opts = &options{
		mockClient: nil,
	}

	for _, inputOption := range inputOptions {
		inputOption(opts)
	}

	if opts.mockClient != nil {
		return &Client{&grpc.ClientConn{}, opts.mockClient}, nil

	}

	c := new(Client)
	sockAddr := fmt.Sprintf("%s:%s", function.Protocol, function.Addr)
	conn, err := grpc.Dial(sockAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to execute grpc.Dial(%q): %w", sockAddr, err)
	}
	c.conn = conn
	c.grpcClt = functionpb.NewUserDefinedFunctionClient(conn)
	return c, nil
}

// CloseConn closes the grpc client connection.
func (c *Client) CloseConn(ctx context.Context) error {
	return c.conn.Close()
}

// IsReady returns true if the grpc connection is ready to use.
func (c *Client) IsReady(ctx context.Context, in *emptypb.Empty) (bool, error) {
	resp, err := c.grpcClt.IsReady(ctx, in)
	if err != nil {
		return false, err
	}
	return resp.GetReady(), nil
}

// DoFn applies a function to each datum element.
func (c *Client) DoFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error) {
	mappedDatumList, err := c.grpcClt.DoFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.DoFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// ReduceFn applies a reduce function to a datum stream.
// TODO: use a channel to accept datumStream
func (c *Client) ReduceFn(ctx context.Context, datumStream []*functionpb.Datum) ([]*functionpb.Datum, error) {
	stream, err := c.grpcClt.ReduceFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.ReduceFn(): %w", err)
	}
	for _, datum := range datumStream {
		if err := stream.Send(datum); err != nil {
			return nil, fmt.Errorf("failed to execute stream.Send(%v): %w", datum, err)
		}
	}
	reducedDatumList, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to execute stream.CloseAndRecv(): %w", err)
	}

	return reducedDatumList.GetElements(), nil
}
