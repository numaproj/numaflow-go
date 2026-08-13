package mapstreamer

import "fmt"

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
	NACK = fmt.Sprintf("%U__NACK__", '\\') // U+005C__NACK__
)

// Message is used to wrap the data return by MapStream functions
type Message struct {
	value       []byte
	keys        []string
	tags        []string
	nackOptions *NackOptions
}

// NewMessage creates a Message with value
func NewMessage(value []byte) Message {
	return Message{value: value}
}

// MessageToDrop creates a Message to be dropped
func MessageToDrop() Message {
	return Message{value: []byte{}, tags: []string{DROP}}
}

// MessageToNack creates a Message that negatively acknowledges the input message,
// requesting redelivery. opts may be nil; when set it carries redelivery options.
func MessageToNack(opts *NackOptions) Message {
	return Message{value: []byte{}, tags: []string{NACK}, nackOptions: opts}
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

// Tags returns message tags
func (m Message) Tags() []string {
	return m.tags
}

// NackOptions returns the message's nack options (nil if not a nack message).
func (m Message) NackOptions() *NackOptions {
	return m.nackOptions
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
