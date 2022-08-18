package client

import (
	"context"
	"fmt"
	"log"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	grpcClt functionpb.UserDefinedFunctionClient
}

func NewClient() *client {
	c := new(client)
	sockAddr := fmt.Sprintf("%s:%s", function.Protocol, function.Addr)
	conn, err := grpc.Dial(sockAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial(%q) failed: %v", sockAddr, err)
	}
	defer conn.Close()
	c.grpcClt = functionpb.NewUserDefinedFunctionClient(conn)
	return c
}

func (c *client) DoFn(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var datum *functionpb.Datum
	mappedDatumList, err := c.grpcClt.DoFn(ctx, datum)
	if err != nil {
		log.Fatalf("client.ReduceFn failed: %v", err)
	}

	// TODO: how do we handle the returned result?
	_ = mappedDatumList
}

func (c *client) ReduceFn(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var datumStream []*functionpb.Datum
	stream, err := c.grpcClt.ReduceFn(ctx)
	if err != nil {
		log.Fatalf("client.ReduceFn failed: %v", err)
	}
	for _, datum := range datumStream {
		if err := stream.Send(datum); err != nil {
			log.Fatalf("client.ReduceFn: stream.Send(%v) failed: %v", datum, err)
		}
	}
	reducedDatumList, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.ReduceFn failed: %v", err)
	}

	// TODO: how do we handle the returned result?
	_ = reducedDatumList

}
