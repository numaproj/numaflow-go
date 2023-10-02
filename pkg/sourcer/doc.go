// Package sourcer implements the server code for User Defined Sourcer in golang.
//
// Example:
/*
    import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcer"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/sourcer"
	)

	// A simple example of a user-defined source.
	// The source maintains an array of messages and implements the `Read()`, `Ack()`, and `Pending()` methods.

	// The `Read(x)` method returns the next x number of messages in the array.
	// The `Ack()` method acknowledges the last batch of messages returned by `Read()`.
	// The `Pending()` method returns 0 to indicate that the simple source always has 0 pending messages.

	func main() {
		simpleSource := NewSimpleSource()
		err := sourcer.NewServer(simpleSource).Start(context.Background())
		if err != nil {
			log.Panic("Failed to start source server : ", err)
		}
	}

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
	func (s *SimpleSource) Pending(_ context.Context) int64 {
		// The simple source always returns zero to indicate no pending records.
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
*/

package sourcer
