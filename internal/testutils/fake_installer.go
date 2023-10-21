package testutils

import (
	"context"
)

type FakeInstaller struct {
	NewVersionError      error
	ExistingVersionError error

	ExistingVersionCalled bool
	NewVersionCalled      bool
}

func (fi *FakeInstaller) NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error {
	fi.NewVersionCalled = true
	return fi.NewVersionError
}

func (fi *FakeInstaller) ExistingVersion(goVersionName string) error {
	fi.ExistingVersionCalled = true
	return fi.ExistingVersionError
}
