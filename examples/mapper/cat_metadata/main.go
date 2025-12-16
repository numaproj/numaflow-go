package main

import (
	"context"
	"fmt"
	"log"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type CatMetadata struct {
}

func (c *CatMetadata) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	um := d.UserMetadata()
	// Print the user metadata for e2e tests in numaflow core
	fmt.Println("Groups at mapper: ", um.Groups())
	if um == nil {
		um = &mapper.UserMetadata{}
	}
	um.CreateGroup("map-group")
	um.AddKV("map-group", "map-key", []byte("map-value"))
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys).WithUserMetadata(um))
}

func main() {
	err := mapper.NewServer(&CatMetadata{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start cat metadata server: ", err)
	}
}
