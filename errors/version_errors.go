// Package errors provides interfaces for custom errors
// across the application.
package errors

import "fmt"

// NoInstalledVersionsError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used when there are not installed Go versions.
type NoInstalledVersionsError struct{}

// Error returns back an error message
func (err *NoInstalledVersionsError) Error() string {
	return "there is no any installed version"
}

// DeleteVersionError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used when a version can't be deleted.
type DeleteVersionError struct {
	// Err is the initial error that was raised.
	Err error

	// Version is the version that was the error was raised for.
	Version string
}

// Error returns back an error message
func (err *DeleteVersionError) Error() string {
	return fmt.Sprintf("an error occurred while deleting %q: %q", err.Version, err.Err.Error())
}

// InstallerNotFoundError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used when an installer couln't be found.
// Could be because of the OS type, the architecture type or even the kind of the file.
type InstallerNotFoundError struct {
	// OS contain the user's Operating System (e.g. `darwin`, `linux`, `windows` etc).
	OS string

	// Arch contains the architecture target (e.g. `386`, `amd64`, `arm64`, `s390x` etc).
	Arch string
}

// Error returns back an error message
func (err *InstallerNotFoundError) Error() string {
	return fmt.Sprintf("installer not found for %q %q", err.OS, err.Arch)
}

// ChecksumNotFoundError is a struct that implements the Error method,
// so can "imitate" and error.
//
// This error should be used when the SHA256 Checksum from API response is empty during downloading the file.
type ChecksumNotFoundError struct {
	// OS contain the user's Operating System (e.g. `darwin`, `linux`, `windows` etc).
	OS string

	// Arch contains the architecture target (e.g. `386`, `amd64`, `arm64`, `s390x` etc).
	Arch string
}

// Error returns back an error message
func (err *ChecksumNotFoundError) Error() string {
	return fmt.Sprintf("checksum not found for %q %q", err.OS, err.Arch)
}
