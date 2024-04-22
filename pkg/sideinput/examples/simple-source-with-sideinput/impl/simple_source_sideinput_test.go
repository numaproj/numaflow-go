package impl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/numaproj/numaflow-go/pkg/sourcer"
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
	offsets []sourcer.Offset
}

func (ar TestAckRequest) Offsets() []sourcer.Offset {
	return ar.offsets
}

func TestNewSimpleSource(t *testing.T) {
	underTest := NewSimpleSource()
	// Prepare a channel to receive messages
	messageCh := make(chan sourcer.Message, 20)
	doneCh := make(chan struct{}) // to know when the go routine is done

	go func() {
		underTest.Read(context.TODO(), TestReadRequest{
			count:   2,
			timeout: time.Second,
		}, messageCh)
		close(doneCh)
	}()
	// We will send messages to GlobalChan
	globalChan <- "test_data_1"
	globalChan <- "test_data_2"
	<-doneCh
	assert.Equal(t, 2, len(messageCh))
	// Try reading 4 more messages
	// Since the previous batch didn't get acked, the data source shouldn't allow us to read more messages
	// We should get 0 messages, meaning the channel only holds the previous 2 messages

	go func() {
		// Read 2 messages
		underTest.Read(context.TODO(), TestReadRequest{
			count:   4,
			timeout: time.Second,
		}, messageCh)
	}()

	assert.Equal(t, 2, len(messageCh))

	// Ack the first batch
	msg1 := <-messageCh
	msg2 := <-messageCh
	assert.Equal(t, "test_data_1", string(msg1.Value()))
	assert.Equal(t, "test_data_2", string(msg2.Value()))

	underTest.Ack(context.TODO(), TestAckRequest{
		offsets: []sourcer.Offset{msg1.Offset(), msg2.Offset()},
	})
	doneCh2 := make(chan struct{}) // to know when the go routine is done

	// Reading 6 more messages by sending in the globalchan
	go func() {
		underTest.Read(context.TODO(), TestReadRequest{
			count:   6,
			timeout: time.Second,
		}, messageCh)
		close(doneCh2)
	}()

	// We will send messages to GlobalChan
	globalChan <- "test_data_1"
	globalChan <- "test_data_2"
	globalChan <- "test_data_3"
	globalChan <- "test_data_4"
	globalChan <- "test_data_5"
	globalChan <- "test_data_6"
	<-doneCh2
	assert.Equal(t, 6, len(messageCh))

	// Ack the second batch
	msg3 := <-messageCh
	msg4 := <-messageCh
	msg5 := <-messageCh
	msg6 := <-messageCh
	msg7 := <-messageCh
	msg8 := <-messageCh
	assert.Equal(t, 0, len(messageCh))
	underTest.Ack(context.TODO(), TestAckRequest{
		offsets: []sourcer.Offset{
			msg3.Offset(), msg4.Offset(),
			msg5.Offset(), msg6.Offset(),
			msg7.Offset(), msg8.Offset(),
		},
	})

}
