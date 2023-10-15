// Package errors provides interfaces for custom errors
// accross the application
package errors

import (
	"fmt"
)

// RequestError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used for failed requests.
type RequestError struct {
	// The status code returned from the failed request.
	StatusCode int
}

// Error returns back an error message
func (err *RequestError) Error() string {
	return fmt.Sprintf("request failed with status %d", err.StatusCode)
}
