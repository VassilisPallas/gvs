package testutils

import (
	"io"
)

type FakeInstallHelper struct {
	CreateExecutableSymlinkError error
	UpdateRecentVersionError     error
	CreateTarFileError           error
	UnzippingError               error
	RenameDirectoryError         error
	RemoveTarFileError           error
	Checksum                     string
}

func (fih FakeInstallHelper) CreateTarFile(content io.ReadCloser) error {
	return fih.CreateTarFileError
}

func (fih FakeInstallHelper) GetTarChecksum() (string, error) {
	return fih.Checksum, nil
}

func (fih FakeInstallHelper) UnzipTarFile() error {
	return fih.UnzippingError
}

func (fih FakeInstallHelper) RenameGoDirectory(goVersionName string) error {
	return fih.RenameDirectoryError
}

func (fih FakeInstallHelper) RemoveTarFile() error {
	return fih.RemoveTarFileError
}

func (fih FakeInstallHelper) CreateExecutableSymlink(goVersionName string) error {
	return fih.CreateExecutableSymlinkError
}

func (fih FakeInstallHelper) UpdateRecentVersion(goVersionName string) error {
	return fih.UpdateRecentVersionError
}
