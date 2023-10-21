// Package errors provides interfaces for custom errors
// across the application
package errors

import "fmt"

// RequestError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used when the checksums do not match with each other.
type ChecksumMisMatchError struct {
	// Checksum contains the SHA256 Checksum from API response.
	Checksum string

	// Hash contains the actual hash from the downloaded file.
	Hash string
}

// Error returns back an error message
func (err *ChecksumMisMatchError) Error() string {
	return fmt.Sprintf("checksums do not match.\nExpected: %q\nGot: %q", err.Checksum, err.Hash)
}
