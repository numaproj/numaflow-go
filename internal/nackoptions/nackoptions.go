package nackoptions

// NackOptions carries per-message redelivery options for a negative
// acknowledgement (nack). All fields are optional; a nil pointer means unset.
type NackOptions struct {
	// Delay is the redelivery delay in milliseconds.
	Delay *uint64
	// MaxDeliveries is the maximum number of redelivery attempts.
	MaxDeliveries *uint32
	// Reason is a human-readable reason for the nack.
	Reason *string
	// Generic values passed as nack options
	NackMap map[string]string
}
