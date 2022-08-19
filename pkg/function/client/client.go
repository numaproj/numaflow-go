package client

import (
	"context"
	"fmt"
	"log"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	grpcClt functionpb.UserDefinedFunctionClient
}

func NewClient() (*Client, error) {
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

func (c *Client) IsReady(ctx context.Context) bool {
	// Notice: This API is EXPERIMENTAL and may be changed or removed in a later
	// release.
	return c.conn.GetState() == connectivity.Ready
}

func (c *Client) CloseConn(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("closing connection timed out.. forcing exit.")
		}
	}()
	return c.conn.Close()
}

func (c *Client) DoFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mappedDatumList, err := c.grpcClt.DoFn(ctx, datum)
	if err != nil {
		return nil, fmt.Errorf("failed to execute c.grpcClt.DoFn(): %w", err)
	}

	return mappedDatumList.Elements, nil
}

// TODO: use a channel to accept datumStream?

func (c *Client) ReduceFn(ctx context.Context, datumStream []*functionpb.Datum) ([]*functionpb.Datum, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stream, err := c.grpcClt.ReduceFn(ctx)
	if err != nil {
		log.Fatalf("failed to execute c.grpcClt.ReduceFn(): %v", err)
	}
	for _, datum := range datumStream {
		if err := stream.Send(datum); err != nil {
			log.Fatalf("failed to execute stream.Send(%v): %v", datum, err)
		}
	}
	reducedDatumList, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to execute stream.CloseAndRecv(): %v", err)
	}

	return reducedDatumList.Elements, nil
}
