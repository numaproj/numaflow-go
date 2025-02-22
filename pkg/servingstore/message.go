package servingstore

// StoredResult is the data stored in the store per origin.
type StoredResult struct {
	origin   string
	payloads []Payload
}

// NewStoredResult creates a new StoreResult from the provided origin and payloads.
func NewStoredResult(origin string, payloads []Payload) StoredResult {
	return StoredResult{origin: origin, payloads: payloads}
}

// Payload is each independent result stored in the Store for the given ID.
type Payload struct {
	value []byte
}

// NewPayload creates a new Payload from the given value.
func NewPayload(value []byte) Payload {
	return Payload{value: value}
}

// StoredResults contains 0, 1, or more StoredResult.
type StoredResults []StoredResult

// StoredResultsBuilder returns an empty instance of StoredResults
func StoredResultsBuilder() StoredResults {
	return StoredResults{}
}

// Append appends a StoredResult
func (r StoredResults) Append(msg StoredResult) StoredResults {
	r = append(r, msg)
	return r
}

// Items returns the StoredResults list
func (r StoredResults) Items() []StoredResult {
	return r
}
