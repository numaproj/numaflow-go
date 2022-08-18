package main

import (
	"context"
	"fmt"
	"log"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/function/client"
)

func main() {
	var ctx = context.Background()
	c, err := client.NewClient()
	defer func() {
		err = c.CloseConn(ctx)
		log.Printf("CloseConn() got an err %v\n", err)
	}()
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("example_client_%d", i)
		list, err := c.DoFn(ctx, &functionpb.Datum{
			Key:   key,
			Value: []byte(`example`),
		})
		if err != nil {
			log.Println(key, err)
		}
		log.Println(key)
		for _, e := range list {
			log.Printf("key:%s, value:%s\n", e.Key, e.Value)
		}
	}

}
