package function

import (
	"time"
)

// MessageT is used to wrap the data return by UDF functions.
// Compared with Message, MessageT contains one more field, the event time, usually extracted from the payload.
type MessageT struct {
	value     []byte
	eventTime time.Time
	keys      []string
	tags      []string
}

// NewMessageT creates a Message with eventTime and value
func NewMessageT(value []byte, eventTime time.Time) MessageT {
	return MessageT{value: value, eventTime: eventTime}
}

// WithKeys is used to assign the keys to messageT
func (m MessageT) WithKeys(keys []string) MessageT {
	m.keys = keys
	return m
}

// WithTags is used to assign the tags to messageT
// tags will be used for conditional forwarding
func (m MessageT) WithTags(tags []string) MessageT {
	m.tags = tags
	return m
}

// EventTime returns message eventTime
func (m MessageT) EventTime() time.Time {
	return m.eventTime
}

// Keys returns message keys
func (m MessageT) Keys() []string {
	return m.keys
}

// Value returns message value
func (m MessageT) Value() []byte {
	return m.value
}

// Tags returns message tags
func (m MessageT) Tags() []string {
	return m.tags
}

// MessageTToDrop creates a MessageT to be dropped
func MessageTToDrop() MessageT {
	return MessageT{eventTime: time.Time{}, value: []byte{}, tags: []string{DROP}}
}

type MessageTs []MessageT

// MessageTsBuilder returns an empty instance of MessageTs
func MessageTsBuilder() MessageTs {
	return MessageTs{}
}

// Append appends a MessageT
func (m MessageTs) Append(msg MessageT) MessageTs {
	m = append(m, msg)
	return m
}

// Items returns the MessageT list
func (m MessageTs) Items() []MessageT {
	return m
}
