package sinker

// Response is the processing result of each message
type Response struct {
	// ID corresponds the ID in the message.
	ID string `json:"id"`
	// Successful or not. If it's false, "err" is expected to be present.
	Success bool `json:"success"`
	// Err represents the error message when "success" is false.
	Err string `json:"err,omitempty"`
	// Fallback is true if the message to be sent to the fallback sink.
	Fallback bool `json:"fallback,omitempty"`
}

type Responses []Response

// ResponsesBuilder returns an empty instance of Responses
func ResponsesBuilder() Responses {
	return Responses{}
}

// Append appends a response
func (r Responses) Append(rep Response) Responses {
	r = append(r, rep)
	return r
}

// Items returns the response list
func (r Responses) Items() []Response {
	return r
}

// ResponseOK creates a successful Response with the given id.
// The Success field is set to true.
func ResponseOK(id string) Response {
	return Response{ID: id, Success: true}
}

// ResponseFailure creates a failed Response with the given id and error message.
// The Success field is set to false and the Err field is set to the provided error message.
func ResponseFailure(id, errMsg string) Response {
	return Response{ID: id, Success: false, Err: errMsg}
}

// ResponseFallback creates a Response with the Fallback field set to true.
// This indicates that the message should be sent to the fallback sink.
func ResponseFallback(id string) Response {
	return Response{ID: id, Fallback: true}
}
