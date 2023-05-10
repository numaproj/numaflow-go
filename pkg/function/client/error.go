package client

import (
	"fmt"
)

// ErrKind represents if the error is retryable
type ErrKind int16

const (
	Retryable    ErrKind = iota // The error is retryable
	NonRetryable                // The error is non-retryable
	Unknown                     // Unknown err kind
)

func (ek ErrKind) String() string {
	switch ek {
	case Retryable:
		return "Retryable"
	case NonRetryable:
		return "NonRetryable"
	case Unknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// UDFError is returned to the main numaflow indicates the status of the error
type UDFError struct {
	ErrKind    ErrKind
	ErrMessage string
}

func (e UDFError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrKind, e.ErrMessage)
}

func (e UDFError) ErrorKind() *ErrKind {
	return &e.ErrKind
}

func (e UDFError) ErrorMessage() string {
	return e.ErrMessage
}

// FromError gets error information from the UDFError
func FromError(err error) (k *ErrKind, s string, ok bool) {
	if err == nil {
		return nil, "", true
	}
	if se, ok := err.(interface {
		ErrorKind() *ErrKind
		ErrorMessage() string
	}); ok {
		return se.ErrorKind(), se.ErrorMessage(), true
	}
	unknown := Unknown
	return &unknown, err.Error(), false
}
