package sideinput

// Message is used to wrap the data return by UserSideInput function.
// It contains the data value for the given side input parameter requested.
type Message struct {
	value []byte
}

// NewMessage creates a new Message with the given value
func NewMessage(value []byte) Message {
	return Message{value: value}
}

// Value returns message value
func (m Message) Value() []byte {
	return m.value
}
