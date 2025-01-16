package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

const DEFAULT_SLEEP_SECONDS = 10

type CatSleep struct {
}

func (e *CatSleep) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {

	sleepSeconds := DEFAULT_SLEEP_SECONDS
	secondsString := os.Getenv("SLEEP_SECONDS")
	if secondsString == "" {
		log.Printf("SLEEP_SECONDS environment variable not set, using default %d seconds\n", DEFAULT_SLEEP_SECONDS)
	} else {
		val, err := strconv.Atoi(secondsString)
		if err != nil {
			log.Printf("SLEEP_SECONDS environment variable %q not an int, using default %d seconds\n", secondsString, DEFAULT_SLEEP_SECONDS)
		} else {
			sleepSeconds = val
		}
	}
	time.Sleep(time.Duration(sleepSeconds) * time.Second)
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
}

func main() {
	err := mapper.NewServer(&CatSleep{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
