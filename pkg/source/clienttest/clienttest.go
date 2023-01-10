package clienttest

import (
	"context"
	"fmt"
	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1/sourcemock"

	"google.golang.org/protobuf/types/known/emptypb"
)

// client contains the grpc client for testing.
type client struct {
	grpcClt sourcepb.UserDefinedSourceTransformerClient
}

// New creates a new mock client object.
func New(c *sourcemock.MockUserDefinedSourceTransformerClient) (*client, error) {
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

// TransformFn applies a transform function to each datum element.
// The datum transformation includes data filtering and event time assignment.
func (c *client) TransformFn(ctx context.Context, datum *sourcepb.Datum) ([]*sourcepb.Datum, error) {
	transformedDatumList, err := c.grpcClt.TransformFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.TransformFn(): %w", err)
	}

	return transformedDatumList.GetElements(), nil
}
