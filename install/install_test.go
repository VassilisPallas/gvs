package install_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/internal/testutils"
)

func TestInstallExistingVersionSucess(t *testing.T) {
	version := "go1.21.0"

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper:    testutils.FakeInstallHelper{},
	}

	err := installer.ExistingVersion(version)

	if err != nil {
		t.Errorf("Error should be nil, instead got '%s'", err.Error())
	}
}

func TestInstallExistingVersionFailedSymlinkCreation(t *testing.T) {
	version := "go1.21.0"
	expectedError := fmt.Errorf("An error occurred while creating the symlink")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			CreateExecutableSymlinkError: expectedError,
		},
	}

	err := installer.ExistingVersion(version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallExistingVersionFailedUpdateVersionFile(t *testing.T) {
	version := "go1.21.0"
	expectedError := fmt.Errorf("An error occurred while updating the version file")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			UpdateRecentVersionError: expectedError,
		},
	}

	err := installer.ExistingVersion(version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionSuccess(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			Checksum: checksum,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err != nil {
		t.Errorf("Error should be nil, instead got '%s'", err.Error())
	}
}

func TestInstallNewVersionFailCreateTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while creating the TAR file")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			CreateTarFileError: expectedError,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailChecksumMissmatch(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			Checksum: "some_other_checksum",
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	expectedError := fmt.Errorf("checksums do not match.\nExpected: %s\nGot: %s", checksum, "some_other_checksum")
	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailUnzipTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while unzipping the TAR file")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			Checksum:       checksum,
			UnzippingError: expectedError,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRenameDirectory(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while renaming the directory")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			Checksum:             checksum,
			RenameDirectoryError: expectedError,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRemoveTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while removing the tart file")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{},
		Helper: testutils.FakeInstallHelper{
			Checksum:           checksum,
			RemoveTarFileError: expectedError,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRequest(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while downloading the file")

	installer := install.Install{
		FileUtils: testutils.FakeFiler{},
		ClientAPI: testutils.FakeGoClientAPI{
			DownloadError: expectedError,
		},
		Helper: testutils.FakeInstallHelper{
			Checksum: checksum,
		},
	}

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
	}
}
