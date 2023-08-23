package sideinput

// SideInputMessage is used to wrap the data return by UserSideInput function.
// It contains the data value for the given side input parameter requested.
type SideInputMessage struct {
	value []byte
}

// NewSideInputMessage creates a new Message with the given value
func NewSideInputMessage(value []byte) SideInputMessage {
	return SideInputMessage{value: value}
}

// Value returns message value
func (m SideInputMessage) Value() []byte {
	return m.value
}
