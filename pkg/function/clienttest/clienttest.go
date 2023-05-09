package clienttest

import (
	"context"
	"fmt"
	"io"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
	"golang.org/x/sync/errgroup"
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

// MapFn applies a function to each datum element.
func (c *client) MapFn(ctx context.Context, datum *functionpb.DatumRequest) ([]*functionpb.DatumResponse, error) {
	mappedDatumList, err := c.grpcClt.MapFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.MapFn(): %w", err)
	}

	return mappedDatumList.GetElements(), nil
}

// MapStreamFn applies a function to each datum element and returns a stream.
func (c *client) MapStreamFn(ctx context.Context, datum *functionpb.DatumRequest, datumCh chan<- *functionpb.DatumResponse) error {
	defer close(datumCh)
	stream, err := c.grpcClt.MapStreamFn(ctx, datum)
	if err != nil {
		return fmt.Errorf("failed to execute c.grpcClt.MapStreamFn(): %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var resp *functionpb.DatumResponse
			resp, err = stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			datumCh <- resp
		}
	}

	return nil
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
