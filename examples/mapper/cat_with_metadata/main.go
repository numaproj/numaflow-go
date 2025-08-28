package main

import (
	"context"
	"log"

	md "github.com/numaproj/numaflow-go/pkg/metadata"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type CatWithMetadata struct {
}

func (c *CatWithMetadata) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	// Start with existing metadata
	metadata := d.Metadata()

	// Initialize metadata if nil
	if metadata == nil {
		metadata = md.New()
	}

	metadata.Set("app", "version", []byte("v1.0"))
	metadata.Set("app", "build", []byte("123"))
	metadata.Set("env", "stage", []byte("production"))

	// Create message with enhanced metadata
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys).WithMetadata(metadata))
}

func main() {
	err := mapper.NewServer(&CatWithMetadata{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start cat with metadata server: ", err)
	}
}
