package sideinput

// Message is used to wrap the data return by UserSideInput function.
// It contains the data value for the given side input parameter requested.
type Message struct {
	value       []byte
	noBroadcast bool
}

// BroadcastMessage creates a new Message with the given value
// This is used to broadcast the message to other side input vertices.
func BroadcastMessage(value []byte) Message {
	return Message{value: value, noBroadcast: false}
}

// NoBroadcastMessage creates a new Message with noBroadcast flag set to true
// This is used to drop the message and not to broadcast it to other side input vertices.
func NoBroadcastMessage() Message {
	return Message{value: []byte{}, noBroadcast: true}
}

// Value returns message value
func (m Message) Value() []byte {
	return m.value
}

// NoBroadcast returns noBroadcast flag
func (m Message) NoBroadcast() bool {
	return m.noBroadcast
}
