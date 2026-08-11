package main

import (
	"context"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type Cat struct {
}

func (c *Cat) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	value := d.Value()
	log.Printf("String val: %s\n", value)

	// check if value is even or odd
	i, err := strconv.Atoi(string(value))
	if err != nil {
		log.Printf("failed to parse value %q as int: %v", value, err)
		return mapper.MessagesBuilder()
	}
	log.Printf("Int: %d\n", i)

	// for even let it go through, stop odd
	if i%2 != 0 {
		log.Printf("ODD\n")
		return mapper.MessagesBuilder().Append(mapper.MessageToDrop())
	}
	log.Printf("Even\n")
	retMes := mapper.MessagesBuilder()
	retMes = retMes.Append(mapper.NewMessage([]byte("Filler-hello")))
	retMes = retMes.Append(mapper.NewMessage(d.Value()).WithKeys(keys))
	retMes = retMes.Append(mapper.NewMessage([]byte("Filler-bye")))
	return retMes
}

func main() {
	err := mapper.NewServer(&Cat{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start cat server: ", err)
	}
}
