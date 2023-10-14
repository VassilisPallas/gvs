// Package errors provides interfaces for custom errors
// accross the application
package errors

import (
	"fmt"
)

// RequestError is a struct that implements the Error method,
// so can "imitate" and error.
//
// RequestError struct accepts the StatusCode field, which can be used on the error message.
//
// This error should be used for failed requests.
type RequestError struct {
	StatusCode int
}

// Error returns back an error message
func (err *RequestError) Error() string {
	return fmt.Sprintf("request failed with status %d", err.StatusCode)
}
