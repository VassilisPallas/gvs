package errors

import "fmt"

type RequestError struct {
	StatusCode int
}

func (err *RequestError) Error() string {
	return fmt.Sprintf("request failed with status %d", err.StatusCode)
}
