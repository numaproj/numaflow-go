package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

const DefaultSleepSeconds = 10

type SlowCat struct {
	sleepSeconds int
}

func (e *SlowCat) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	time.Sleep(time.Duration(e.sleepSeconds) * time.Second)
	return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
}

func main() {
	err := mapper.NewServer(&SlowCat{sleepSeconds: getSleepSeconds()}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}

func getSleepSeconds() int {

	// sleep time is configured according to environment variable (or default if not configured)

	sleepSeconds := DefaultSleepSeconds
	secondsString := os.Getenv("SLEEP_SECONDS")
	if secondsString == "" {
		log.Printf("SLEEP_SECONDS environment variable not set, using default %d seconds\n", DefaultSleepSeconds)
	} else {
		val, err := strconv.Atoi(secondsString)
		if err != nil {
			log.Printf("SLEEP_SECONDS environment variable %q not an int, using default %d seconds\n", secondsString, DefaultSleepSeconds)
		} else {
			sleepSeconds = val
			log.Printf("Using SLEEP_SECONDS value %d\n", sleepSeconds)
		}
	}
	return sleepSeconds
}
