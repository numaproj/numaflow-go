package source

import (
	"fmt"
	"time"
)

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
	ALL  = fmt.Sprintf("%U__ALL__", '\\')  // U+005C__ALL__
)

// Message is used to wrap the data return by transform functions
type Message struct {
	EventTime time.Time
	Key       string
	Value     []byte
}

// MessageToDrop creates a Message to be dropped
func MessageToDrop() Message {
	return Message{Key: DROP, Value: []byte{}}
}

// MessageToAll creates a Message that will forward to all
func MessageToAll(et time.Time, value []byte) Message {
	return Message{EventTime: et, Key: ALL, Value: value}
}

// MessageTo creates a Message that will forward to specified "to"
func MessageTo(et time.Time, to string, value []byte) Message {
	return Message{EventTime: et, Key: to, Value: value}
}

type Messages []Message

// MessagesBuilder returns an empty instance of Messages
func MessagesBuilder() Messages {
	return Messages{}
}

// Append appends a Message
func (m Messages) Append(msg Message) Messages {
	m = append(m, msg)
	return m
}

// Items returns the message list
func (m Messages) Items() []Message {
	return m
}
