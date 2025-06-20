package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type Cat struct {
}

func (c *Cat) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
}

func main() {
	err := mapper.NewServer(&Cat{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start cat server: ", err)
	}
}
