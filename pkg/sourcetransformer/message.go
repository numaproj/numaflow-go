package sourcetransformer

import (
	"fmt"
	"time"
)

var DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
var NACK = fmt.Sprintf("%U__NACK__", '\\') // U+005C__NACK__

// Message is used to wrap the data return by SourceTransformer functions.
// Compared with Message of other UDFs, source transformer Message contains one more field,
// the event time, usually extracted from the payload.
type Message struct {
	value        []byte
	eventTime    time.Time
	keys         []string
	tags         []string
	userMetadata *UserMetadata
	nackOptions  *NackOptions
}

// NewMessage creates a Message with eventTime and value
func NewMessage(value []byte, eventTime time.Time) Message {
	return Message{value: value, eventTime: eventTime}
}

// WithKeys is used to assign the keys to message
func (m Message) WithKeys(keys []string) Message {
	m.keys = keys
	return m
}

// WithTags is used to assign the tags to message
// tags will be used for conditional forwarding
func (m Message) WithTags(tags []string) Message {
	m.tags = tags
	return m
}

// WithUserMetadata is used to assign the user metadata to the message
func (m Message) WithUserMetadata(userMetadata *UserMetadata) Message {
	m.userMetadata = userMetadata
	return m
}

// EventTime returns message eventTime
func (m Message) EventTime() time.Time {
	return m.eventTime
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

// MessageToDrop creates a Message to be dropped with eventTime.
// eventTime is required because, even though a message is dropped, it is still considered as being processed,
// hence the watermark should be updated accordingly using the provided event time.
func MessageToDrop(eventTime time.Time) Message {
	return Message{eventTime: eventTime, value: []byte{}, tags: []string{DROP}}
}

// MessageToNack creates a Message that negatively acknowledges the input message,
// requesting redelivery. eventTime is required (the watermark still advances).
// opts may be nil; when set it carries redelivery options.
func MessageToNack(eventTime time.Time, opts *NackOptions) Message {
	return Message{eventTime: eventTime, value: []byte{}, tags: []string{NACK}, nackOptions: opts}
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

// Items returns the Message list
func (m Messages) Items() []Message {
	return m
}

// UserMetadata returns message user metadata
func (m Message) UserMetadata() *UserMetadata {
	return m.userMetadata
}
