package sinker

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
