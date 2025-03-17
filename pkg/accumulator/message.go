package accumulator

import (
	"fmt"
	"time"
)

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
)

// Message is used to wrap the data return by Map functions
type Message struct {
	value     []byte
	keys      []string
	tags      []string
	eventTime time.Time
	watermark time.Time
	id        string
	headers   map[string]string
}

// MessageFromDatum creates a Message from the input datum. It's advised to use the same input datum for creating the
// message, only use a custom implementation of the Datum if you know what you are doing.
func MessageFromDatum(datum Datum) Message {
	return Message{
		value:     datum.Value(),
		eventTime: datum.EventTime(),
		watermark: datum.Watermark(),
		id:        datum.ID(),
		keys:      datum.Keys(),
		headers:   datum.Headers(),
	}
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

// WithValue is used to assign the value to the message
func (m Message) WithValue(value []byte) Message {
	m.value = value
	return m
}

// WithID is used to assign the ID to the message, ID will be used for deduplication, it's better to use the ID from the
// input datum. Only use a custom ID if you know what you are doing.
func (m Message) WithID(id string) Message {
	m.id = id
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
