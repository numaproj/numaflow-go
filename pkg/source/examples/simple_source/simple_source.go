package main

import (
	"context"
	"strconv"
	"time"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/model"
)

type SimpleSource struct {
	size     int64
	messages []string
	readIdx  int64
	toAckSet map[int64]struct{}
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
	}
}

func (s *SimpleSource) Pending(_ context.Context) uint64 {
	// The simple source always returns 0 to indicate no pending records.
	return 0
}

func (s *SimpleSource) Read(_ context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- model.Message) {
	// Handle the timeout specification in the read request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*readRequest.TimeOut())
	defer cancel()

	// If we have already read all the data,
	// Or we have un-acked data, we return.
	if s.readIdx >= s.size || len(s.toAckSet) > 0 {
		return
	}

	counter := readRequest.Count()
	// Read the data from the source and send the data to the message channel.
	select {
	case <-ctx.Done():
		if !toAckSetMatchReadCount(s.toAckSet, readRequest.Count()-counter) {
			panic("toAckSet does not match read count")
		}
		// If the context is done, the read request is timed out.
		return
	default:
		if counter <= 0 {
			if !toAckSetMatchReadCount(s.toAckSet, readRequest.Count()) {
				panic("toAckSet does not match read count")
			}
			return
		}
		// Otherwise, we read the data from the source and send the data to the message channel.
		offsetValue := serializeOffset(s.readIdx)
		messageCh <- model.NewMessage(
			[]byte(s.messages[s.readIdx]),
			model.NewOffset(offsetValue, "0"),
			time.Now())
		// Mark the offset as to be acked.
		s.toAckSet[s.readIdx] = struct{}{}
		s.readIdx++
		counter--
		time.Sleep(1 * time.Millisecond)
	}
}

func (s *SimpleSource) Ack(_ context.Context, request sourcesdk.AckRequest) {
	// TODO - either ack all or none
	offsetsToAck := request.Offsets()
	if len(offsetsToAck) != len(s.toAckSet) {
		panic("offsets to ack do not match the toAckSet")
	}
	for _, offset := range offsetsToAck {
		delete(s.toAckSet, deserializeOffset(offset.Value()))
	}
}

func toAckSetMatchReadCount(toAckSet map[int64]struct{}, readCount uint64) bool {
	return len(toAckSet) == int(readCount)
}

func serializeOffset(idx int64) []byte {
	return []byte(strconv.FormatInt(idx, 10))
}

func deserializeOffset(offset []byte) int64 {
	idx, _ := strconv.ParseInt(string(offset), 10, 64)
	return idx
}
