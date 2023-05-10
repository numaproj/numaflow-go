package client

import (
	"fmt"
)

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
		return "Non-retryable"
	case Unknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

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
	return &unknown, "", false
}
