package mapper

import (
	"fmt"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
	md "github.com/numaproj/numaflow-go/pkg/metadata"
)

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
)

// Message is used to wrap the data return by Map functions
type Message struct {
	value    []byte
	keys     []string
	tags     []string
	metadata *md.Metadata
}

// NewMessage creates a Message with value
func NewMessage(value []byte) Message {
	return Message{value: value}
}

// MessageToDrop creates a Message to be dropped
func MessageToDrop() Message {
	return Message{value: []byte{}, tags: []string{DROP}}
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

// WithMetadata is used to assign the metadata to the message
func (m Message) WithMetadata(metadata *md.Metadata) Message {
	m.metadata = metadata
	return m
}

// Metadata returns message metadata
func (m Message) Metadata() *md.Metadata {
	return m.metadata
}

// toProto converts md.Metadata to protobuf MessageMetadata
func toProto(m *md.Metadata) *mappb.MessageMetadata {
	if m == nil {
		return nil
	}
	return m.MessageMetadata
}

// fromProto creates mdpb.Metadata from protobuf MessageMetadata

func fromProto(proto *mappb.MessageMetadata) *md.Metadata {
	if proto == nil {
		return nil
	}
	return &md.Metadata{MessageMetadata: proto}
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
