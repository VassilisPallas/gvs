package testutils

import (
	"context"
)

type FakeInstaller struct {
	NewVersionError      error
	ExistingVersionError error
}

func (fi FakeInstaller) NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error {
	return fi.NewVersionError
}

func (fi FakeInstaller) ExistingVersion(goVersionName string) error {
	return fi.ExistingVersionError
}
