package function

import "fmt"

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
)

// Message is used to wrap the data return by UDF functions
type Message struct {
	tags  []string
	keys  []string
	value []byte
}

// NewMessage creates a Message with value
func NewMessage(value []byte) Message {
	return Message{value: value}
}

// MessageToDrop creates a Message to be dropped
func MessageToDrop() Message {
	return Message{tags: []string{DROP}, value: []byte{}}
}

// WithKeys is used to assign the keys to the message
func (m Message) WithKeys(keys []string) Message {
	m.keys = keys
	return m
}

// WithTags is used to assign the tags to the message
// tags will be used for conditional forwarding
func (m Message) WithTags(tags []string) Message {
	m.tags = tags
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
