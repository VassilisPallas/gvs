package install_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/internal/testutils"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type nopCloserWriter struct {
	writerError error

	io.Reader
}

func (nopCloserWriter) Close() error { return nil }

func (ncw nopCloserWriter) WriteTo(w io.Writer) (n int64, err error) {
	return n, ncw.writerError
}

func TestCreateTarFile(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	fileContent := "foo"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	err := helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

	defer os.Remove(fileLocation)

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	content, err := os.ReadFile(fileLocation)
	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	if string(content) != fileContent {
		t.Errorf("content should be '%s', instead got '%s'", fileContent, content)
	}
}

func TestCreateTarFileToPathThatDoesNotExist(t *testing.T) {
	fileLocation := "/some_other_dst/some_file.tar.gz"
	fileContent := "foo"
	expectedError := "open /some_other_dst/some_file.tar.gz: no such file or directory"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	err := helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}
}

func TestCreateTarFileCopyFailed(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	fileContent := "foo"
	expectedError := "some error while copying"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	err := helper.CreateTarFile(nopCloserWriter{writerError: fmt.Errorf(expectedError), Reader: bytes.NewBufferString(fileContent)})

	defer os.Remove(fileLocation)

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}
}

func TestGetTarChecksum(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	fileContent := "foo"
	expectedHash := "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	// Create file for the test
	helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})
	defer os.Remove(fileLocation)

	hash, err := helper.GetTarChecksum()

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	if hash != expectedHash {
		t.Errorf("hash should be '%s', instead got '%s'", expectedHash, hash)
	}
}

func TestGetTarChecksumFileToPathThatDoesNotExist(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	expectedError := "open /tmp/some_file.tar.gz: no such file or directory"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	hash, err := helper.GetTarChecksum()

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}

	if hash != "" {
		t.Errorf("hash should be an empty string, instead got '%s'", hash)
	}
}

func TestRenameGoDirectory(t *testing.T) {
	dirLocation := "/tmp/some_dir"
	goVersion := "go1.21.0"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			VersionDir: dirLocation,
		},
	}

	os.MkdirAll(fmt.Sprintf("%s/%s", dirLocation, "go"), 0755)
	defer os.RemoveAll(dirLocation)

	err := helper.RenameGoDirectory(goVersion)

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	expectedDirectory := fmt.Sprintf("%s/%s", dirLocation, goVersion)
	if _, err := os.Stat(expectedDirectory); os.IsNotExist(err) {
		t.Errorf("%s does not exist", expectedDirectory)
	}
}

func TestRemoveTarFile(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	fileContent := "foo"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	// Create file for the test
	helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

	err := helper.RemoveTarFile()

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	if _, err := os.Stat(fileLocation); !os.IsNotExist(err) {
		t.Errorf("%s exists while it should be deleted", fileLocation)
	}
}

func TestRemoveTarFileToPathThatDoesNotExist(t *testing.T) {
	fileLocation := "/tmp/some_file.tar.gz"
	expectedError := "remove /tmp/some_file.tar.gz: no such file or directory"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			TarFile: fileLocation,
		},
	}

	err := helper.RemoveTarFile()

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}
}

func TestUpdateRecentVersion(t *testing.T) {
	fileLocation := "/tmp/CURRENT"
	goVersion := "go1.21.0"

	helper := install.Helper{
		FileUtils: testutils.FakeFiler{
			CurrentVersionFile: fileLocation,
		},
	}

	err := helper.UpdateRecentVersion(goVersion)
	defer os.Remove(fileLocation)

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	content, err := os.ReadFile(fileLocation)
	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	if string(content) != goVersion {
		t.Errorf("content should be '%s', instead got '%s'", goVersion, content)
	}
}
