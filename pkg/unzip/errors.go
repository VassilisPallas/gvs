// Package unzip provides functions for extracting
// tar and zip files.
package unzip

import (
	"errors"
	"fmt"
	"io"
)

// UnzipError is a struct that implements the Error method,
// so can "imitate" and error.
type UnzipError struct {
	// This is the initial error that is wrapped on the custom error
	err error
}

// Error returns back an error message
func (err *UnzipError) Error() string {
	return fmt.Sprintf("unzip failed with error: %s", err.err.Error())
}

// isEOF returns is the error is EOF.
//
// isEOF is used to stop the iterations until we reach the end of the read
func isEOF(err error) bool {
	return errors.Is(err, io.EOF)
}
