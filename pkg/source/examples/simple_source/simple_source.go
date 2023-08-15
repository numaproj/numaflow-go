package main

import (
	"context"
	"strconv"
	"sync"
	"time"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/model"
)

// SimpleSource is a simple source implementation.
type SimpleSource struct {
	size     int64
	messages []string
	readIdx  int64
	toAckSet map[int64]struct{}
	lock     *sync.Mutex
}

// NewSimpleSource creates a SimpleSource by populating the messages array with random strings
func NewSimpleSource(size int64) *SimpleSource {
	messages := make([]string, size)
	for i := int64(0); i < size; i++ {
		messages[i] = strconv.FormatInt(i, 10)
	}
	return &SimpleSource{
		size:     size,
		messages: messages,
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
func (s *SimpleSource) Read(_ context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- model.Message) {
	// Handle the timeout specification in the read request.
	ctx, cancel := context.WithTimeout(context.Background(), readRequest.TimeOut())
	defer cancel()

	// If we have un-acked data, or if we have read all the data,
	// we return without reading any new data.
	if len(s.toAckSet) > 0 || s.readIdx >= s.size {
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
			messageCh <- model.NewMessage(
				[]byte(s.messages[s.readIdx]),
				model.NewOffset(offsetValue, "0"),
				time.Now())
			// Mark the offset as to be acked, and increment the read index.
			s.toAckSet[s.readIdx] = struct{}{}
			s.readIdx++
			s.lock.Unlock()
			if s.readIdx >= s.size {
				// If we have read all the data, we return.
				return
			}
		}
	}
}

// Ack acknowledges the data from the source.
// If the offsets in the request do not match the offsets to be acked, the function panics.
// Note:
// If the method is about to panic at any point of the execution,
// before it panics, it ensures that all the offsets in the request are NOT acked.
// This is to maintain the contract of the Ack method - either acknowledge ALL or NONE.
func (s *SimpleSource) Ack(_ context.Context, request sourcesdk.AckRequest) {
	offsetsToAck := request.Offsets()
	if len(offsetsToAck) != len(s.toAckSet) {
		panic("offsets to ack do not match the toAckSet")
	}

	copyOfOriginalToAckSet := make(map[int64]struct{})
	for k, v := range s.toAckSet {
		copyOfOriginalToAckSet[k] = v
	}

	for _, offset := range offsetsToAck {
		dk := deserializeOffset(offset.Value())
		if _, ok := s.toAckSet[dk]; !ok {
			// One of the input offsets is not in the toAckSet.
			// Revoke all the previous acknowledgements and panic.
			s.toAckSet = copyOfOriginalToAckSet
			panic("offsets to ack do not match the toAckSet")
		} else {
			// Remove the offset from the toAckSet.
			delete(s.toAckSet, dk)
		}
	}
}

func serializeOffset(idx int64) []byte {
	return []byte(strconv.FormatInt(idx, 10))
}

func deserializeOffset(offset []byte) int64 {
	idx, _ := strconv.ParseInt(string(offset), 10, 64)
	return idx
}
