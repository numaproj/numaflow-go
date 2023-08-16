package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/numaproj/numaflow-go/pkg/source"
)

// Notes: the unit test cases below demonstrate the basic contract between Read and Ack functions of a data source.
// These test cases should be applicable to all data sources, not just SimpleSource.

type TestReadRequest struct {
	count   uint64
	timeout time.Duration
}

func (rr TestReadRequest) Count() uint64 {
	return rr.count
}

func (rr TestReadRequest) TimeOut() time.Duration {
	return rr.timeout
}

type TestAckRequest struct {
	offsets []source.Offset
}

func (ar TestAckRequest) Offsets() []source.Offset {
	return ar.offsets
}

func Test_SimpleSource(t *testing.T) {
	underTest := NewSimpleSource()
	// Prepare a channel to receive messages
	messageCh := make(chan source.Message, 20)

	// Read 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 2, len(messageCh))

	// Try reading 4 more messages
	// Since the previous batch didn't get acked, the data source shouldn't allow us to read more messages
	// We should get 0 messages, meaning the channel only holds the previous 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   4,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 2, len(messageCh))

	// Ack the first batch
	msg1 := <-messageCh
	msg2 := <-messageCh
	underTest.Ack(nil, TestAckRequest{
		offsets: []source.Offset{msg1.Offset(), msg2.Offset()},
	})

	// Try reading 6 more messages
	// Since the previous batch got acked, the data source should allow us to read more messages
	// We should get 6 messages
	underTest.Read(nil, TestReadRequest{
		count:   6,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 6, len(messageCh))

	// Ack the second batch
	msg3 := <-messageCh
	msg4 := <-messageCh
	msg5 := <-messageCh
	msg6 := <-messageCh
	msg7 := <-messageCh
	msg8 := <-messageCh
	assert.Equal(t, 0, len(messageCh))
	underTest.Ack(nil, TestAckRequest{
		offsets: []source.Offset{
			msg3.Offset(), msg4.Offset(),
			msg5.Offset(), msg6.Offset(),
			msg7.Offset(), msg8.Offset(),
		},
	})
}
