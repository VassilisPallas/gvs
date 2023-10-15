package files_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/pkg/unzip"
)

func createFileHelper(cliWriter io.Writer, logWriter io.WriteCloser, fs testutils.FakeFileSystem, unzipper unzip.Unzipper) *files.Helper {
	logger := logger.New(cliWriter, logWriter)
	return files.New(fs, unzipper, logger)
}

func TestCreateTarFile(t *testing.T) {
	testCases := []struct {
		testTitle   string
		createError error
		copyError   error
	}{
		{
			testTitle:   "creates tar file",
			createError: nil,
			copyError:   nil,
		},
		{
			testTitle:   "create file fails",
			createError: errors.New("an error occurred while creating the file"),
			copyError:   nil,
		},
		{
			testTitle:   "copy content fails",
			createError: nil,
			copyError:   errors.New("an error occurred while creating the file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			var file *os.File

			if tc.createError == nil {
				var creationError error
				file, creationError = os.Create("/tmp/someFile.tar.gz") // create file for the testing purposes
				if creationError != nil {
					t.Fatal("could not create the file to run the test")
				}
			}

			// remove file after the test if it was created
			defer func() {
				if tc.createError == nil {
					os.Remove("/tmp/someFile.tar.gz")
				}
			}()

			fs := testutils.FakeFileSystem{
				CreateMockFile: file,
				CreateError:    tc.createError,
				CopyError:      tc.copyError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{})

			fileContent := "foo"
			content := io.NopCloser(bytes.NewBufferString(fileContent))
			err := fileHelper.CreateTarFile(content)

			if tc.copyError == nil && tc.createError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.createError != nil && err.Error() != tc.createError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.createError.Error(), err.Error())
				return
			}

			if tc.copyError != nil && err.Error() != tc.copyError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.copyError.Error(), err.Error())
				return
			}
		})
	}
}

func TestGetTarChecksum(t *testing.T) {
	testCases := []struct {
		testTitle string
		openError error
		copyError error
	}{
		{
			testTitle: "returns tar checksum",
			openError: nil,
			copyError: nil,
		},
		{
			testTitle: "open file fails",
			openError: errors.New("an error occurred while opening the fails"),
			copyError: nil,
		},
		{
			testTitle: "copy of hash fails",
			openError: nil,
			copyError: errors.New("an error occurred while copying the contents"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			var file *os.File

			if tc.openError == nil {
				var creationError error
				file, creationError = os.Create("/tmp/someFile.tar.gz") // create file for the testing purposes
				if creationError != nil {
					t.Fatal("could not create the file to run the test")
				}
			}

			// remove file after the test if it was created
			defer func() {
				if tc.openError == nil {
					os.Remove("/tmp/someFile.tar.gz")
				}
			}()

			fs := testutils.FakeFileSystem{
				OpenMockFile: file,
				OpenError:    tc.openError,
				CopyError:    tc.copyError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{})

			hash, err := fileHelper.GetTarChecksum()

			if tc.copyError == nil && tc.openError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.openError != nil && err.Error() != tc.openError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.openError.Error(), err.Error())
				return
			}

			if tc.copyError != nil && err.Error() != tc.copyError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.copyError.Error(), err.Error())
				return
			}

			expectedHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
			if tc.copyError == nil && tc.openError == nil && hash != expectedHash {
				t.Errorf("hash should be %q, instead got %q", expectedHash, hash)
				return
			}
		})
	}
}

func TestUnzipTarFile(t *testing.T) {
	testCases := []struct {
		testTitle string
	}{}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {

		})
	}
}

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"
// 	"testing"

// 	"github.com/VassilisPallas/gvs/api_client"
// 	"github.com/VassilisPallas/gvs/files"
// 	"github.com/VassilisPallas/gvs/internal/testutils"
// )

// type nopCloser struct {
// 	io.Reader
// }

// func (nopCloser) Close() error { return nil }

// type nopCloserWriter struct {
// 	writerError error

// 	io.Reader
// }

// func (nopCloserWriter) Close() error { return nil }

// func (ncw nopCloserWriter) WriteTo(w io.Writer) (n int64, err error) {
// 	return n, ncw.writerError
// }

// func TestCreateTarFile(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	fileContent := "foo"

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

// 	defer os.Remove(fileLocation)

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	content, err := os.ReadFile(fileLocation)
// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	if string(content) != fileContent {
// 		t.Errorf("content should be '%s', instead got '%s'", fileContent, content)
// 	}
// }

// func TestCreateTarFileToPathThatDoesNotExist(t *testing.T) {
// 	fileLocation := "/some_other_dst/some_file.tar.gz"
// 	fileContent := "foo"
// 	expectedError := fmt.Errorf("open /some_other_dst/some_file.tar.gz: no such file or directory")

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestCreateTarFileCopyFailed(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	fileContent := "foo"
// 	expectedError := fmt.Errorf("some error while copying")

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.CreateTarFile(nopCloserWriter{writerError: expectedError, Reader: bytes.NewBufferString(fileContent)})

// 	defer os.Remove(fileLocation)

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestGetTarChecksum(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	fileContent := "foo"
// 	expectedHash := "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	// Create file for the test
// 	helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})
// 	defer os.Remove(fileLocation)

// 	hash, err := helper.GetTarChecksum()

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	if hash != expectedHash {
// 		t.Errorf("hash should be '%s', instead got '%s'", expectedHash, hash)
// 	}
// }

// func TestGetTarChecksumFileToPathThatDoesNotExist(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	expectedError := fmt.Errorf("open /tmp/some_file.tar.gz: no such file or directory")

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	hash, err := helper.GetTarChecksum()

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}

// 	if hash != "" {
// 		t.Errorf("hash should be an empty string, instead got '%s'", hash)
// 	}
// }

// func TestRenameGoDirectory(t *testing.T) {
// 	dirLocation := "/tmp/some_dir"
// 	goVersion := "go1.21.0"

// 	fileUtils := testutils.FakeFiler{
// 		VersionDir: dirLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	os.MkdirAll(fmt.Sprintf("%s/%s", dirLocation, "go"), 0755)
// 	defer os.RemoveAll(dirLocation)

// 	err := helper.RenameGoDirectory(goVersion)

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	expectedDirectory := fmt.Sprintf("%s/%s", dirLocation, goVersion)
// 	if _, err := os.Stat(expectedDirectory); os.IsNotExist(err) {
// 		t.Errorf("%s does not exist", expectedDirectory)
// 	}
// }

// func TestRemoveTarFile(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	fileContent := "foo"

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	// Create file for the test
// 	helper.CreateTarFile(nopCloser{bytes.NewBufferString(fileContent)})

// 	err := helper.RemoveTarFile()

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	if _, err := os.Stat(fileLocation); !os.IsNotExist(err) {
// 		t.Errorf("%s exists while it should be deleted", fileLocation)
// 	}
// }

// func TestRemoveTarFileToPathThatDoesNotExist(t *testing.T) {
// 	fileLocation := "/tmp/some_file.tar.gz"
// 	expectedError := fmt.Errorf("remove /tmp/some_file.tar.gz: no such file or directory")

// 	fileUtils := testutils.FakeFiler{
// 		TarFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.RemoveTarFile()

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestUpdateRecentVersion(t *testing.T) {
// 	fileLocation := "/tmp/CURRENT"
// 	goVersion := "go1.21.0"

// 	fileUtils := testutils.FakeFiler{
// 		CurrentVersionFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.UpdateRecentVersion(goVersion)
// 	defer os.Remove(fileLocation)

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	content, err := os.ReadFile(fileLocation)
// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	if string(content) != goVersion {
// 		t.Errorf("content should be '%s', instead got '%s'", goVersion, content)
// 	}
// }

// func TestStoreVersionsResponse(t *testing.T) {
// 	appDir := "/tmp"
// 	versionResponseFile := "goVersions.json"

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.StoreVersionsResponse([]byte("Some content"))

// 	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}
// }

// func TestStoreVersionsResponseShouldNotFailWithEmptyContent(t *testing.T) {
// 	appDir := "/tmp"
// 	versionResponseFile := "goVersions.json"

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)
// 	err := helper.StoreVersionsResponse([]byte(""))

// 	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}
// }

// func TestStoreVersionsResponseFileToPathThatDoesNotExist(t *testing.T) {
// 	appDir := "/some_other_dst"
// 	versionResponseFile := "goVersions.json"
// 	expectedError := fmt.Errorf("open /some_other_dst/goVersions.json: no such file or directory")

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)

// 	err := helper.StoreVersionsResponse([]byte("Some content"))

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestGetCachedResponse(t *testing.T) {
// 	appDir := "/tmp"
// 	versionResponseFile := "goVersions.json"

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)
// 	responseVersionsMap := []map[string]interface{}{
// 		{
// 			"version": "go1.21.0",
// 			"stable":  true,
// 			"files": []any{
// 				map[string]any{
// 					"arch":     string("arm64"),
// 					"filename": string("go1.21.0.linux-arm64.tar.gz"),
// 					"kind":     string("archive"),
// 					"os":       string("linux"),
// 					"sha256":   string("818d46ede85682dd551ad378ef37a4d247006f12ec59b5b755601d2ce114369a"),
// 					"size":     float64(9.6962473e+07),
// 					"version":  string("go1.21.0"),
// 				},
// 				map[string]any{
// 					"arch":     string("amd64"),
// 					"filename": string("go1.21.0.darwin-amd64.pkg"),
// 					"kind":     string("archive"),
// 					"os":       string("darwin"),
// 					"sha256":   string("725de310e4cba0121d6337053b2cfc3fe47da4a3d50726582731cb1d2a70137e"),
// 					"size":     float64(6.714125e+07),
// 					"version":  string("go1.21.0"),
// 				},
// 			},
// 		}}

// 	responseBody, _ := json.Marshal(responseVersionsMap)

// 	helper.StoreVersionsResponse([]byte(responseBody))

// 	var responseVersions []api_client.VersionInfo
// 	err := helper.GetCachedResponse(&responseVersions)

// 	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}

// 	if responseVersions == nil {
// 		t.Errorf("responseVersions should not be nil")
// 	}
// }

// func TestGetCachedResponseFileToPathThatDoesNotExist(t *testing.T) {
// 	appDir := "/some_other_dst"
// 	versionResponseFile := "goVersions.json"
// 	expectedError := fmt.Errorf("open /some_other_dst/goVersions.json: no such file or directory")

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)

// 	var responseVersions []api_client.VersionInfo
// 	err := helper.GetCachedResponse(&responseVersions)

// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestGetCachedResponseUnmarshalFailed(t *testing.T) {
// 	appDir := "/tmp"
// 	versionResponseFile := "goVersions.json"

// 	fileUtils := testutils.FakeFiler{
// 		AppDir:                  appDir,
// 		VersionResponseFileName: versionResponseFile,
// 	}
// 	helper := files.New(fileUtils)

// 	helper.StoreVersionsResponse([]byte("{foo: bar}")) // force syntax error to response body

// 	var responseVersions []api_client.VersionInfo
// 	err := helper.GetCachedResponse(&responseVersions)

// 	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

// 	expectedError := fmt.Errorf("invalid character 'f' looking for beginning of object key string")
// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }

// func TestGetRecentVersion(t *testing.T) {
// 	fileLocation := "/tmp/CURRENT"
// 	goVersion := "go1.21.0"

// 	fileUtils := testutils.FakeFiler{
// 		CurrentVersionFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	file, _ := os.Create(fileLocation)
// 	io.WriteString(file, goVersion)

// 	content := helper.GetRecentVersion()

// 	if content != goVersion {
// 		t.Errorf("content is not correct. Got '%s', instead of '%s'", content, goVersion)
// 	}
// }

// func TestGetRecentVersionFileToPathThatDoesNotExist(t *testing.T) {
// 	fileLocation := "/some_other_dst/CURRENT"

// 	fileUtils := testutils.FakeFiler{
// 		CurrentVersionFile: fileLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	content := helper.GetRecentVersion()

// 	if content != "" {
// 		t.Errorf("content should be empty string, instead it got '%s'", content)
// 	}
// }

// func TestDirectoryExists(t *testing.T) {
// 	dirLocation := "/tmp"
// 	goVersion := "go1.21.0"

// 	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

// 	fileUtils := testutils.FakeFiler{
// 		VersionDir: dirLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	os.MkdirAll(path, 0755)

// 	exists := helper.DirectoryExists(goVersion)

// 	defer os.RemoveAll(path)

// 	if !exists {
// 		t.Errorf("%s should exist", path)
// 	}
// }

// func TestDirectoryExistsDirectoryNotFound(t *testing.T) {
// 	dirLocation := "/tmp"
// 	goVersion := "go1.21.0"

// 	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

// 	fileUtils := testutils.FakeFiler{
// 		VersionDir: dirLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	exists := helper.DirectoryExists(goVersion)

// 	if exists {
// 		t.Errorf("%s should not exist", path)
// 	}
// }

// func TestDeleteDirectory(t *testing.T) {
// 	dirLocation := "/tmp"
// 	goVersion := "go1.21.0"

// 	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

// 	fileUtils := testutils.FakeFiler{
// 		VersionDir: dirLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	os.MkdirAll(path, 0755)

// 	err := helper.DeleteDirectory(goVersion)

// 	defer os.RemoveAll(path)

// 	if err != nil {
// 		t.Errorf("error should be nil, instead got '%s'", err.Error())
// 	}
// }

// func TestDeleteDirectoryDeleteFails(t *testing.T) {
// 	dirLocation := "/tmp"
// 	goVersion := "." // This is forcing os.RemoveAll to fail

// 	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

// 	fileUtils := testutils.FakeFiler{
// 		VersionDir: dirLocation,
// 	}
// 	helper := files.New(fileUtils)

// 	os.MkdirAll(path, 0755)

// 	err := helper.DeleteDirectory(goVersion)

// 	defer os.RemoveAll(path)

// 	expectedError := fmt.Errorf("RemoveAll /tmp/.: invalid argument")
// 	if err.Error() != expectedError.Error() {
// 		t.Errorf("error should be '%s', instead got '%s'", expectedError.Error(), err.Error())
// 	}
// }
