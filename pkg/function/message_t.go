package function

import (
	"time"
)

// MessageT is used to wrap the data return by UDF functions.
// Compared with Message, MessageT contains one more field, the event time, usually extracted from the payload.
type MessageT struct {
	EventTime time.Time
	Key       string
	Value     []byte
}

// MessageTToDrop creates a MessageT to be dropped
func MessageTToDrop() MessageT {
	return MessageT{EventTime: time.Time{}, Key: DROP, Value: []byte{}}
}

// MessageTToAll creates a MessageT that will forward to all
func MessageTToAll(eventTime time.Time, value []byte) MessageT {
	return MessageT{EventTime: eventTime, Key: ALL, Value: value}
}

// MessageTTo creates a MessageT that will forward to specified "to"
func MessageTTo(eventTime time.Time, to string, value []byte) MessageT {
	return MessageT{EventTime: eventTime, Key: to, Value: value}
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
