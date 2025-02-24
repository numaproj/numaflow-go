package servingstore

// StoredResult is the data stored in the store per origin.
type StoredResult struct {
	id       string
	payloads []Payload
}

// NewStoredResult creates a new StoreResult from the provided origin and payloads.
func NewStoredResult(id string, payloads []Payload) StoredResult {
	return StoredResult{id: id, payloads: payloads}
}

// Payload is each independent result stored in the Store for the given ID.
type Payload struct {
	origin string
	value  []byte
}

// NewPayload creates a new Payload from the given value.
func NewPayload(origin string, value []byte) Payload {
	return Payload{origin: origin, value: value}
}

// Value returns the value of the Payload.
func (p *Payload) Value() []byte {
	return p.value
}

// Origin returns the origin name.
func (p *Payload) Origin() string {
	return p.origin
}
