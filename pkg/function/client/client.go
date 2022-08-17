package function

import (
	"context"
	"log"
	"net"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/grpc"
)

func NewClient() *client {
	c := new(client)

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial(protocol, addr)
	}

	conn, err := grpc.Dial(sockAddr, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c.grpcClt = functionpb.NewUserDefinedFunctionClient(conn)

	return c
}

type client struct {
	grpcClt functionpb.UserDefinedFunctionClient
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
