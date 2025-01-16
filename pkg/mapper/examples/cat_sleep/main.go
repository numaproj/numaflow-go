package main

import (
	"context"
	//"fmt"
	"log"
	"time"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type CatSleep struct {
}

func (e *CatSleep) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	time.Sleep(10 * time.Second)
	//fmt.Println("adding a log line")
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
}

func main() {
	err := mapper.NewServer(&CatSleep{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
