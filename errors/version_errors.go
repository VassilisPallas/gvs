package errors

import "fmt"

type NoInstalledVersionsError struct{}

func (err *NoInstalledVersionsError) Error() string {
	return "there is no any installed version"
}

type DeleteVersionError struct {
	Err     error
	Version string
}

func (err *DeleteVersionError) Error() string {
	return fmt.Sprintf("an error occurred while deleting %q: %q", err.Version, err.Err.Error())
}

type InstalledNotFoundError struct {
	OS   string
	Arch string
}

func (err *InstalledNotFoundError) Error() string {
	return fmt.Sprintf("installer not found for %q %q", err.OS, err.Arch)
}

type ChecksumNotFoundError struct {
	OS   string
	Arch string
}

func (err *ChecksumNotFoundError) Error() string {
	return fmt.Sprintf("checksum not found for %q %q", err.OS, err.Arch)
}
