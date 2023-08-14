package model

import "time"

// Message is used to wrap the data return by UDSource
type Message struct {
	value     []byte
	offset    Offset
	eventTime time.Time
	keys      []string
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

// Keys returns message keys
func (m Message) Keys() []string {
	return m.keys
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
	partitionId string
}

// NewOffset creates an Offset with value and partition id
func NewOffset(value []byte, partitionId string) Offset {
	return Offset{value: value, partitionId: partitionId}
}

// Value returns value of the offset
func (o Offset) Value() []byte {
	return o.value
}

// PartitionId returns partition id of the offset
func (o Offset) PartitionId() string {
	return o.partitionId
}
