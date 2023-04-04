package function

import "fmt"

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
	ALL  = fmt.Sprintf("%U__ALL__", '\\')  // U+005C__ALL__
)

// Message is used to wrap the data return by UDF functions
type Message struct {
	keys  []string
	value []byte
}

// Keys returns message keys
func (m *Message) Keys() []string {
	return m.keys
}

// Value returns message value
func (m *Message) Value() []byte {
	return m.value
}

// MessageToDrop creates a Message to be dropped
func MessageToDrop() Message {
	return Message{keys: []string{DROP}, value: []byte{}}
}

// MessageToAll creates a Message that will forward to all
func MessageToAll(value []byte) Message {
	return Message{keys: []string{ALL}, value: value}
}

// MessageTo creates a Message that will forward to specified "to"
func MessageTo(to []string, value []byte) Message {
	return Message{keys: to, value: value}
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
