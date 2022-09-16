package clienttest

import (
	"context"
	"fmt"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1/sinkmock"
	"google.golang.org/protobuf/types/known/emptypb"
)

// client contains the grpc client for testing.
type client struct {
	grpcClt sinkpb.UserDefinedSinkClient
}

// New creates a new mock client object.
func New(c *sinkmock.MockUserDefinedSinkClient) (*client, error) {
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

// SinkFn applies a function to a list of datum elements.
func (c *client) SinkFn(ctx context.Context, datumList []*sinkpb.Datum) ([]*sinkpb.Response, error) {
	responseList, err := c.grpcClt.SinkFn(ctx, &sinkpb.DatumList{Elements: datumList})
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.SinkFn(): %w", err)
	}

	return responseList.GetResponses(), nil
}
