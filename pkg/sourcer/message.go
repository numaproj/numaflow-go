package sourcer

import (
	"os"
	"strconv"
	"time"
)

// create default partition id from the environment variable "NUMAFLOW_REPLICA"
var defaultPartitionId, _ = strconv.Atoi(os.Getenv("NUMAFLOW_REPLICA"))

// Message is used to wrap the data return by UDSource
type Message struct {
	value     []byte
	offset    Offset
	eventTime time.Time
	keys      []string
	headers   map[string]string
}

// NewMessage creates a Message with value
func NewMessage(value []byte, offset Offset, eventTime time.Time) Message {
	return Message{value: value, offset: offset, eventTime: eventTime}
}

// WithKeys is used to assign the keys to the message
func (m Message) WithKeys(keys []string) Message {
	m.keys = keys
	return m
}

// WithHeaders is used to assign the headers to the message
func (m Message) WithHeaders(headers map[string]string) Message {
	m.headers = headers
	return m
}

// Keys returns message keys
func (m Message) Keys() []string {
	return m.keys
}

// Headers returns message headers
func (m Message) Headers() map[string]string {
	return m.headers
}

// Value returns message value
func (m Message) Value() []byte {
	return m.value
}

// Offset returns message offset
func (m Message) Offset() Offset {
	return m.offset
}

// EventTime returns message event time
func (m Message) EventTime() time.Time {
	return m.eventTime
}

type Offset struct {
	value       []byte
	partitionId int32
}

// NewOffset creates an Offset with value and partition id
func NewOffset(value []byte, partitionId int32) Offset {
	return Offset{value: value, partitionId: partitionId}
}

// NewOffsetWithDefaultPartitionId creates an Offset with value and default partition id. This
// function can be used if you use DefaultPartitions() to implement the Sourcer interface.
// For most cases, this function can be used as long as the source does not have a concept of partitions.
// If you need to implement a custom partition, use `NewOffset`.
func NewOffsetWithDefaultPartitionId(value []byte) Offset {
	return Offset{value: value, partitionId: DefaultPartitions()[0]}
}

// DefaultPartitions returns default partitions for the source.
// It can be used in the Partitions() function of the Sourcer interface only
// if the source doesn't have partitions. DefaultPartition will be the pod replica
// index of the source.
func DefaultPartitions() []int32 {
	return []int32{int32(defaultPartitionId)}
}

// Value returns value of the offset
func (o Offset) Value() []byte {
	return o.value
}

// PartitionId returns partition id of the offset
func (o Offset) PartitionId() int32 {
	return o.partitionId
}
