package impl

import (
	"context"
	"strconv"
	"sync"
	"time"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/sourcer"
)

// SimpleSource is a simple source implementation.
type SimpleSource struct {
	readIdx  int64
	toAckSet map[int64]struct{}
	lock     *sync.Mutex
}

func NewSimpleSource() *SimpleSource {
	return &SimpleSource{
		readIdx:  0,
		toAckSet: make(map[int64]struct{}),
		lock:     new(sync.Mutex),
	}
}

// Pending returns the number of pending records.
func (s *SimpleSource) Pending(_ context.Context) uint64 {
	// The simple source always returns 0 to indicate no pending records.
	return 0
}

// Read reads messages from the source and sends the messages to the message channel.
// If the read request is timed out, the function returns without reading new data.
// Right after reading a message, the function marks the offset as to be acked.
func (s *SimpleSource) Read(_ context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- sourcesdk.Message) {
	// Handle the timeout specification in the read request.
	ctx, cancel := context.WithTimeout(context.Background(), readRequest.TimeOut())
	defer cancel()

	// If we have un-acked data, we return without reading any new data.
	// TODO - we should have a better way to handle this.
	// In real world, there can be the case when the numa container restarted before acknowledging a batch,
	// leaving the toAckSet not empty on the UDSource container side.
	// In this case, for the next batch read, we should read the data from the last acked offset instead of returning.
	// Our built-in Kafka source follows this logic.
	if len(s.toAckSet) > 0 {
		return
	}

	// Read the data from the source and send the data to the message channel.
	for i := 0; uint64(i) < readRequest.Count(); i++ {
		select {
		case <-ctx.Done():
			// If the context is done, the read request is timed out.
			return
		default:
			s.lock.Lock()
			// Otherwise, we read the data from the source and send the data to the message channel.
			offsetValue := serializeOffset(s.readIdx)
			messageCh <- sourcesdk.NewMessage(
				[]byte(strconv.FormatInt(s.readIdx, 10)),
				sourcesdk.NewOffset(offsetValue, "0"),
				time.Now())
			// Mark the offset as to be acked, and increment the read index.
			s.toAckSet[s.readIdx] = struct{}{}
			s.readIdx++
			s.lock.Unlock()
		}
	}
}

// Ack acknowledges the data from the source.
func (s *SimpleSource) Ack(_ context.Context, request sourcesdk.AckRequest) {
	for _, offset := range request.Offsets() {
		delete(s.toAckSet, deserializeOffset(offset.Value()))
	}
}

func serializeOffset(idx int64) []byte {
	return []byte(strconv.FormatInt(idx, 10))
}

func deserializeOffset(offset []byte) int64 {
	idx, _ := strconv.ParseInt(string(offset), 10, 64)
	return idx
}
