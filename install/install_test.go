package install_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/google/go-cmp/cmp"
)

func TestInstallExistingVersionSuccess(t *testing.T) {
	version := "go1.21.0"

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.ExistingVersion(version)

	if err != nil {
		t.Errorf("Error should be nil, instead got %q", err.Error())
	}
}

func TestInstallExistingVersionSuccessLogs(t *testing.T) {
	printer := &testutils.FakeStdout{}

	version := "go1.21.0"

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(printer, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.ExistingVersion(version)

	if err != nil {
		t.Errorf("Error should be nil, instead got %q", err.Error())
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"Installing version...\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestInstallExistingVersionFailedSymlinkCreation(t *testing.T) {
	version := "go1.21.0"
	expectedError := fmt.Errorf("An error occurred while creating the symlink")

	fileHelpers := &testutils.FakeFilesHelper{
		CreateExecutableSymlinkError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.ExistingVersion(version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallExistingVersionFailedUpdateVersionFile(t *testing.T) {
	version := "go1.21.0"
	expectedError := fmt.Errorf("An error occurred while updating the version file")

	fileHelpers := &testutils.FakeFilesHelper{
		UpdateRecentVersionError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.ExistingVersion(version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionSuccess(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum: checksum,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err != nil {
		t.Errorf("Error should be nil, instead got %q", err.Error())
	}
}

func TestInstallNewVersionSuccessLogs(t *testing.T) {
	printer := &testutils.FakeStdout{}

	version := "go1.21.0"
	checksum := "some_checksum"

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum: checksum,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(printer, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err != nil {
		t.Errorf("Error should be nil, instead got %q", err.Error())
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"Downloading...\n",
		"Compare Checksums...\n",
		"Unzipping...\n",
		"Installing version...\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestInstallNewVersionFailCreateTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while creating the TAR file")

	fileHelpers := &testutils.FakeFilesHelper{
		CreateTarFileError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailGetFileChecksum(t *testing.T) {
	expectedError := fmt.Errorf("some error")
	checksum := ""

	version := "go1.21.0"

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum:            checksum,
		GetTarChecksumError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if !fileHelpers.RemoveTarFileCalled {
		t.Errorf("RemoveTarFileCalled has not been called")
	}

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailChecksumMissmatch(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum: "some_other_checksum",
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	expectedError := fmt.Errorf("checksums do not match.\nExpected: %q\nGot: %q", checksum, "some_other_checksum")
	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailUnzipTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while unzipping the TAR file")

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum:       checksum,
		UnzippingError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRenameDirectory(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while renaming the directory")

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum:             checksum,
		RenameDirectoryError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRemoveTarFile(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while removing the tart file")

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum:           checksum,
		RemoveTarFileError: expectedError,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFailRequest(t *testing.T) {
	version := "go1.21.0"
	checksum := "some_checksum"

	expectedError := fmt.Errorf("An error occurred while downloading the file")

	fileHelpers := &testutils.FakeFilesHelper{
		Checksum: checksum,
	}
	clientAPI := testutils.FakeGoClientAPI{
		DownloadError: expectedError,
	}
	logger := logger.New(&testutils.FakeStdout{}, nil)

	installer := install.New(fileHelpers, clientAPI, logger)

	err := installer.NewVersion(context.Background(), "some_file_name", checksum, version)

	if err.Error() != expectedError.Error() {
		t.Errorf("Error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}
