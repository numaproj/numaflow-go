package clienttest

import (
	"context"
	"fmt"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
	"google.golang.org/protobuf/types/known/emptypb"
)

// client contains the grpc client for testing.
type client struct {
	grpcClt functionpb.UserDefinedFunctionClient
}

// New creates a new mock client object.
func New(c *funcmock.MockUserDefinedFunctionClient) (*client, error) {
	return &client{c}, nil
}

// CloseConn closes the grpc client connection.
func (c *client) CloseConn(ctx context.Context) error {
	return nil
}

// IsReady returns true if the grpc connection is ready to use.
func (c *client) IsReady(ctx context.Context, in *emptypb.Empty) (bool, error) {
	resp, err := c.grpcClt.IsReady(ctx, in)
	if err != nil {
		return false, err
	}
	return resp.GetReady(), nil
}

// MapFn applies a function to each datum element without modifying event time.
func (c *client) MapFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error) {
	mappedDatumList, err := c.grpcClt.MapFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.MapFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// MapTFn applies a function to each datum element.
// In addition to map function, MapTFn also supports assigning a new event time to datum.
func (c *client) MapTFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error) {
	mappedDatumList, err := c.grpcClt.MapTFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.MapTFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// ReduceFn applies a reduce function to a datum stream.
func (c *client) ReduceFn(ctx context.Context, datumStreamCh <-chan *functionpb.Datum) ([]*functionpb.Datum, error) {
	stream, err := c.grpcClt.ReduceFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.ReduceFn(): %w", err)
	}
	for datum := range datumStreamCh {
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
