package errors

import "fmt"

type ChecksumMisMatchError struct {
	Checksum string
	Hash     string
}

func (err *ChecksumMisMatchError) Error() string {
	return fmt.Sprintf("checksums do not match.\nExpected: %q\nGot: %q", err.Checksum, err.Hash)
}
