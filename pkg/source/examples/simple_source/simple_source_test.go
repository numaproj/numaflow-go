package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/numaproj/numaflow-go/pkg/source/model"
)

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
	offsets []model.Offset
}

func (ar TestAckRequest) Offsets() []model.Offset {
	return ar.offsets
}

func Test_SimpleSource_HappyPath(t *testing.T) {
	// Create a new UDSource with 10 messages
	underTest := NewSimpleSource(10)
	// Prepare a channel to receive messages
	// The channel size is 20 to make sure it can hold all messages
	messageCh := make(chan model.Message, 20)

	// Read 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 2, len(messageCh))

	// Try reading 2 more messages
	// Since the previous batch didn't get acked, the data source shouldn't allow us to read more messages
	// We should get 0 messages, meaning the channel only holds the previous 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 2, len(messageCh))

	// Ack the first batch
	msg1 := <-messageCh
	msg2 := <-messageCh
	assert.Equal(t, 0, len(messageCh))
	underTest.Ack(nil, TestAckRequest{
		offsets: []model.Offset{msg1.Offset(), msg2.Offset()},
	})

	// Try reading 2 more messages
	// Since the previous batch got acked, the data source should allow us to read more messages
	// We should get 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 2, len(messageCh))

	// Ack the second batch
	msg3 := <-messageCh
	msg4 := <-messageCh
	assert.Equal(t, 0, len(messageCh))
	underTest.Ack(nil, TestAckRequest{
		offsets: []model.Offset{msg3.Offset(), msg4.Offset()},
	})

	// At this point, the data source should have 6 messages left
	// Try reading 10 messages, it should return 6 messages
	underTest.Read(nil, TestReadRequest{
		count:   10,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 6, len(messageCh))

	// Ack the last batch
	msg5 := <-messageCh
	msg6 := <-messageCh
	msg7 := <-messageCh
	msg8 := <-messageCh
	msg9 := <-messageCh
	msg10 := <-messageCh
	assert.Equal(t, 0, len(messageCh))
	underTest.Ack(nil, TestAckRequest{
		offsets: []model.Offset{msg5.Offset(), msg6.Offset(), msg7.Offset(), msg8.Offset(), msg9.Offset(), msg10.Offset()},
	})

	// At this point, the data source should have 0 messages left
	// Try reading another batch should return 0 messages
	underTest.Read(nil, TestReadRequest{
		count:   10,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 0, len(messageCh))
}

func Test_SimpleSource_NotSoHappyPath(t *testing.T) {
	// Create a new UDSource with 10 messages
	underTest := NewSimpleSource(10)
	// Prepare a channel to receive messages
	// The channel size is 20 to make sure it can hold all messages
	messageCh := make(chan model.Message, 20)

	// Read 2 messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)

	// Try Acknowledging the batch, but passing in wrong offsets
	// The AckRequest contains 3 offsets, which is more than the number of messages we read
	// The Ack function should ack none of the messages and panic
	batch1Msg1 := <-messageCh
	batch1Msg2 := <-messageCh
	randomMsg := model.NewMessage([]byte("random"), model.NewOffset([]byte("random"), "0"), time.Now())
	assertPanic(t, func() {
		underTest.Ack(nil, TestAckRequest{
			offsets: []model.Offset{batch1Msg1.Offset(), batch1Msg2.Offset(), randomMsg.Offset()},
		})
	})
	// Assert that we can't read more messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 0, len(messageCh))

	// Try Acknowledging the batch, but passing in wrong offsets
	// The AckRequest contains 2 offsets, but one of them doesn't match the message we read
	// The Ack function should ack none of the messages and panic
	assertPanic(t, func() {
		underTest.Ack(nil, TestAckRequest{
			offsets: []model.Offset{batch1Msg1.Offset(), randomMsg.Offset()},
		})
	})
	// Assert that we can't read more messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 0, len(messageCh))

	// Try Acknowledging the batch, but passing in wrong offsets
	// The AckRequest contains 1 offset, it matches one of the messages we read
	// The Ack function should ack none of the messages and panic, because we can't ack a subset of the batch
	assertPanic(t, func() {
		underTest.Ack(nil, TestAckRequest{
			offsets: []model.Offset{batch1Msg1.Offset()},
		})
	})
	// Assert that we can't read more messages
	underTest.Read(nil, TestReadRequest{
		count:   2,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 0, len(messageCh))

	// Finally, Ack the batch correctly
	underTest.Ack(nil, TestAckRequest{
		offsets: []model.Offset{batch1Msg1.Offset(), batch1Msg2.Offset()},
	})
	// Assert that we can read more messages
	underTest.Read(nil, TestReadRequest{
		count:   3,
		timeout: time.Second,
	}, messageCh)
	assert.Equal(t, 3, len(messageCh))
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
