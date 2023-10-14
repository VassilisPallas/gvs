// Package errors provides interfaces for custom errors
// accross the application
package errors

import "fmt"

// RequestError is a struct that implements the Error method,
// so can "imitate" and error.
//
// RequestError struct accepts two fields, the Checksum field, which that contains the SHA256 Checksum from API response
// and the Hash, which contains the actual hash from the downloaded file.
//
// This error should be used when the checksums do not match with each other.
type ChecksumMisMatchError struct {
	Checksum string
	Hash     string
}

// Error returns back an error message
func (err *ChecksumMisMatchError) Error() string {
	return fmt.Sprintf("checksums do not match.\nExpected: %q\nGot: %q", err.Checksum, err.Hash)
}
