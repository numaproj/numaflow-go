package mapper

import "fmt"

var (
	DROP = fmt.Sprintf("%U__DROP__", '\\') // U+005C__DROP__
)

// ===========================================================================================
// Common structures used in map functions
// ===========================================================================================

// Message is used to wrap the data return by Map functions
type Message struct {
	value []byte
	keys  []string
	tags  []string
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

// ===========================================================================================
// Utility structures for unary map use case
// ===========================================================================================

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

// ===========================================================================================
// Utility structures for batch map mode
// ===========================================================================================

// batchResponse is used to wrap the data return by batch map function along
// with the ID of the corresponding request
type batchResponse struct {
	id       string
	messages []Message
}

// Id returns request ID for the given list of responses
func (m batchResponse) Id() string {
	return m.id
}

// Append appends a Message to the messages list of a batchResponse
// object and then returns the updated object.
func (m batchResponse) Append(msg Message) batchResponse {
	m.messages = append(m.messages, msg)
	return m
}

// Items returns the message list for a batchResponse
func (m batchResponse) Items() []Message {
	return m.messages
}

// BatchResponses is a list of batchResponse which signify the consolidated
// results for a batch of input messages.
type BatchResponses []batchResponse

// NewBatchResponse is a utility function used to create a new batchResponse object
// Specifying an id is a mandatory requirement, as it is required to reference the
// responses back to a request.
func NewBatchResponse(id string) batchResponse {
	return batchResponse{
		id:       id,
		messages: MessagesBuilder(),
	}
}
func BatchResponsesBuilder() BatchResponses {
	return BatchResponses{}
}

// Append appends a Message
func (m BatchResponses) Append(msg batchResponse) BatchResponses {
	m = append(m, msg)
	return m
}

// Items returns the message list
func (m BatchResponses) Items() []batchResponse {
	return m
}
