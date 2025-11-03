package sinker

import (
	"context"
	"time"
)

// Message is used to wrap the data returned by Sink
type Message struct {
	value        []byte
	keys         []string
	userMetadata *UserMetadata
}

// NewMessage creates a Message with value
func NewMessage(value []byte) Message {
	return Message{value: value}
}

// WithKeys is used to assign the keys to the message
func (m Message) WithKeys(keys []string) Message {
	m.keys = keys
	return m
}

// WithUserMetadata is used to assign the user metadata to the message
func (m Message) WithUserMetadata(userMetadata *UserMetadata) Message {
	m.userMetadata = userMetadata
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

// UserMetadata returns message user metadata
func (m Message) UserMetadata() *UserMetadata {
	return m.userMetadata
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

// Datum is the interface of incoming message payload for sink function.
type Datum interface {
	// Keys returns the keys of the message.
	Keys() []string
	// Value returns the payload of the message.
	Value() []byte
	// EventTime returns the event time of the message.
	EventTime() time.Time
	// Watermark returns the watermark of the message.
	Watermark() time.Time
	// ID returns the ID of the message.
	ID() string
	// Headers returns the headers of the message.
	Headers() map[string]string
	// UserMetadata returns the user metadata of the message.
	UserMetadata() *UserMetadata
	// SystemMetadata returns the system metadata of the message.
	SystemMetadata() *SystemMetadata
}

// Sinker is the interface of sink function implementation.
type Sinker interface {
	// Sink is the function to process a list of incoming messages
	Sink(ctx context.Context, datumStreamCh <-chan Datum) Responses
}

// SinkerFunc is utility type used to convert a Sink function to a Sinker.
type SinkerFunc func(ctx context.Context, datumStreamCh <-chan Datum) Responses

// Sink implements the function of sink function.
func (sf SinkerFunc) Sink(ctx context.Context, datumStreamCh <-chan Datum) Responses {
	return sf(ctx, datumStreamCh)
}
