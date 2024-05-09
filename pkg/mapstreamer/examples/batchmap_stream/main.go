package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/mapstreamer"
)

// Get hash of byte array and serialized bytes
func hash(d mapstreamer.Datum) (uint64, []byte) {

	msg := d.Value()
	summer := fnv.New64a()
	summer.Write(msg)
	val := summer.Sum64()

	asStrBytes := []byte(strconv.FormatUint(val, 10))
	return val, asStrBytes
}

// BatchMap is a MapStreamer that split the input message into multiple messages and stream them.
type BatchMap struct {
}

func (f *BatchMap) MapStream(ctx context.Context, keys []string, d mapstreamer.Datum, messageCh chan<- mapstreamer.Message) {

	defer close(messageCh)

	val, b := hash(d)

	if val%2 == 0 {
		messageCh <- mapstreamer.NewMessage(b).WithKeys([]string{"even"}).WithTags([]string{"even-tag"})
	} else {
		messageCh <- mapstreamer.NewMessage(b).WithKeys([]string{"odd"}).WithTags([]string{"odd-tag"})
	}
}

func (f *BatchMap) MapStreamBatch(ctx context.Context, datumCh <-chan mapstreamer.Datum, messageCh chan<- mapstreamer.Message) {
	count := 0

	grouped := make([]string, 0)

	// Go through each message and figure out if the hash is odd or even and output as individual streamed message
	for d := range datumCh {
		count += 1

		val, b := hash(d)

		if val%2 == 0 {
			messageCh <- mapstreamer.NewMessage(b).WithKeys([]string{"even"}).WithTags([]string{"even-tag"})
		} else {
			messageCh <- mapstreamer.NewMessage(b).WithKeys([]string{"odd"}).WithTags([]string{"odd-tag"})
		}

		// Also track hashes of each
		grouped = append(grouped, fmt.Sprint(val))
	}

	// Demonstrate can also send separate message that is aggregate of the batch
	asJson, _ := json.Marshal(grouped)
	messageCh <- mapstreamer.NewMessage(asJson).WithTags([]string{"grouped"})
}

func main() {
	err := mapstreamer.NewServer(&BatchMap{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map stream function server: ", err)
	}
}
