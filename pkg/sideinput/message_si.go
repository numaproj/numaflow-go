package sideinput

// MessageSI is used to wrap the data return by UserSideInput function.
// It contains the data value for the given side input parameter requested.
type MessageSI struct {
	value []byte
}

// NewMessageSI creates a new Message with the given value
func NewMessageSI(value []byte) MessageSI {
	return MessageSI{value: value}
}

// Value returns message value
func (m MessageSI) Value() []byte {
	return m.value
}
