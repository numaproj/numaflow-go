package function

import (
	"time"
)

// MessageT is used to wrap the data return by UDF functions.
// Compared with Message, MessageT contains one more field, the event time, usually extracted from the payload.
type MessageT struct {
	eventTime time.Time
	key       []string
	value     []byte
}

// EventTime returns message eventTime
func (m *MessageT) EventTime() time.Time {
	return m.eventTime
}

// Key returns message key
func (m *MessageT) Key() []string {
	return m.key
}

// Value returns message value
func (m *MessageT) Value() []byte {
	return m.value
}

// MessageTToDrop creates a MessageT to be dropped
func MessageTToDrop() MessageT {
	return MessageT{eventTime: time.Time{}, key: []string{DROP}, value: []byte{}}
}

// MessageTToAll creates a MessageT that will forward to all
func MessageTToAll(eventTime time.Time, value []byte) MessageT {
	return MessageT{eventTime: eventTime, key: []string{ALL}, value: value}
}

// MessageTTo creates a MessageT that will forward to specified "to"
func MessageTTo(eventTime time.Time, to []string, value []byte) MessageT {
	return MessageT{eventTime: eventTime, key: to, value: value}
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
