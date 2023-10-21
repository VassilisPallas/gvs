package files_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	ioFS "io/fs"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/clock"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/pkg/unzip"
)

func createFileHelper(cliWriter io.Writer, logWriter io.WriteCloser, fs testutils.FakeFileSystem, unzipper unzip.Unzipper, clock clock.Clock) *files.Helper {
	logger := logger.New(cliWriter, logWriter)
	return files.New(fs, clock, unzipper, logger)
}

func getDirEntries(fakeEntries []testutils.FakeDirEntry) []ioFS.DirEntry {
	entries := make([]ioFS.DirEntry, len(fakeEntries))
	for i := range fakeEntries {
		entries[i] = fakeEntries[i]
	}

	return entries
}

func getGoDirEntry(DirEntryInfoError error) []testutils.FakeDirEntry {
	return []testutils.FakeDirEntry{
		{
			DirEntryName:  "go",
			DirEntryIsDir: true,
			DireEntryType: 0,
			DirEntryInfo: testutils.FakeFileInfo{
				FileName:    "go",
				FileSize:    1000,
				FileMode:    0,
				FileModTime: time.Time{},
				FileIsDir:   true,
			},
			DirEntryInfoError: DirEntryInfoError,
		},
	}
}

func getEmptyDirEntries() []testutils.FakeDirEntry {
	return []testutils.FakeDirEntry{}
}

func getGoVersionBinEntries() []testutils.FakeDirEntry {
	return []testutils.FakeDirEntry{
		{
			DirEntryName:  "go",
			DirEntryIsDir: false,
			DireEntryType: 0,
			DirEntryInfo: testutils.FakeFileInfo{
				FileName:    "go",
				FileSize:    1000,
				FileMode:    0,
				FileModTime: time.Time{},
				FileIsDir:   false,
			},
			DirEntryInfoError: nil,
		},
		{
			DirEntryName:  "gofmt",
			DirEntryIsDir: false,
			DireEntryType: 0,
			DirEntryInfo: testutils.FakeFileInfo{
				FileName:    "gofmt",
				FileSize:    1000,
				FileMode:    0,
				FileModTime: time.Time{},
				FileIsDir:   false,
			},
			DirEntryInfoError: nil,
		},
	}
}

func TestCreateTarFile(t *testing.T) {
	testCases := []struct {
		testTitle     string
		createError   error
		copyError     error
		expectedError error
	}{
		{
			testTitle:     "should create the tar file",
			createError:   nil,
			copyError:     nil,
			expectedError: nil,
		},
		{
			testTitle:     "should fail to create the tar file with an error when creating occurs",
			createError:   errors.New("an error occurred while creating the file"),
			copyError:     nil,
			expectedError: errors.New("an error occurred while creating the file"),
		},
		{
			testTitle:     "should fail to create the tar file with an error when copy the contents occurs",
			createError:   nil,
			copyError:     errors.New("an error occurred while creating the file"),
			expectedError: errors.New("an error occurred while creating the file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			var file *os.File

			fs := testutils.FakeFileSystem{
				CreateMockFile: file,
				CreateError:    tc.createError,
				CopyError:      tc.copyError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})

			fileContent := "foo"
			content := io.NopCloser(bytes.NewBufferString(fileContent))
			err := fileHelper.CreateTarFile(content)

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}
		})
	}
}

func TestGetTarChecksum(t *testing.T) {
	testCases := []struct {
		testTitle     string
		openError     error
		copyError     error
		expectedError error
		expectedHash  string
	}{
		{
			testTitle:     "should return tar checksum",
			openError:     nil,
			copyError:     nil,
			expectedError: nil,
			expectedHash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			testTitle:     "should return an empty string when opening the tar file fails",
			openError:     errors.New("an error occurred while opening the fails"),
			copyError:     nil,
			expectedError: errors.New("an error occurred while opening the fails"),
			expectedHash:  "",
		},
		{
			testTitle:     "copy of hash fails",
			openError:     nil,
			copyError:     errors.New("should return an empty string when copying the content of the checksum fails"),
			expectedError: errors.New("should return an empty string when copying the content of the checksum fails"),
			expectedHash:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			var file *os.File

			fs := testutils.FakeFileSystem{
				OpenMockFile: file,
				OpenError:    tc.openError,
				CopyError:    tc.copyError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})

			hash, err := fileHelper.GetTarChecksum()

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}

			if hash != tc.expectedHash {
				t.Errorf("hash should be %q, instead got %q", tc.expectedHash, hash)
				return
			}
		})
	}
}

func TestUnzipTarFile(t *testing.T) {
	testCases := []struct {
		testTitle  string
		unzipError error
	}{
		{
			testTitle:  "should successfully unzip the tar file",
			unzipError: nil,
		},
		{
			testTitle:  "should return an error when unzipping the tar file fails",
			unzipError: errors.New("an error occurred while unzipping"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{}
			unzipper := testutils.FakeUnzipper{ExtractTarSourceError: tc.unzipError}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, unzipper, testutils.FakeClock{})

			err := fileHelper.UnzipTarFile()
			if tc.unzipError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.unzipError != nil && err.Error() != tc.unzipError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.unzipError.Error(), err.Error())
				return
			}
		})
	}
}

func TestRenameGoDirectory(t *testing.T) {
	testCases := []struct {
		testTitle       string
		readDirResponse []testutils.FakeDirEntry
		readDirError    error
		renameError     error
		expectedError   error
	}{
		{
			testTitle:       "should rename the directory",
			readDirResponse: getGoDirEntry(nil),
			readDirError:    nil,
			renameError:     nil,
			expectedError:   nil,
		},
		{
			testTitle:       "should return an error when unable to get the latest created directory",
			readDirResponse: getEmptyDirEntries(),
			readDirError:    errors.New("an error happened while reading the directories"),
			renameError:     nil,
			expectedError:   errors.New("an error happened while reading the directories"),
		},
		{
			testTitle:       "should return an error when renaming the directory fails",
			readDirResponse: getGoDirEntry(nil),
			readDirError:    nil,
			renameError:     errors.New("an error happened while renaming the directory"),
			expectedError:   errors.New("an error happened while renaming the directory"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:             "/tmp",
				ReadDirMockResponse: getDirEntries(tc.readDirResponse),
				ReadDirError:        tc.readDirError,
				RenameError:         tc.renameError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			err := fileHelper.RenameGoDirectory("go1.21.0")

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}
		})
	}
}

func TestRemoveTarFile(t *testing.T) {
	testCases := []struct {
		testTitle   string
		removeError error
	}{
		{
			testTitle:   "should remove the tar file",
			removeError: nil,
		},
		{
			testTitle:   "should return an error when deleting the tar file fails",
			removeError: errors.New("an error happened while removing the tar file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:     "/tmp",
				RemoveError: tc.removeError,
			}

			unzipper := testutils.FakeUnzipper{}
			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, unzipper, testutils.FakeClock{})
			err := fileHelper.RemoveTarFile()

			if tc.removeError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.removeError != nil && err.Error() != tc.removeError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.removeError.Error(), err.Error())
				return
			}
		})
	}
}

func TestCreateExecutableSymlink(t *testing.T) {
	testCases := []struct {
		testTitle       string
		readDirResponse []testutils.FakeDirEntry
		readDirError    error
		removeError     error
		symlinkError    error
		chmodError      error
		expectedError   error
	}{
		{
			testTitle:       "should create the symlink",
			readDirResponse: getGoVersionBinEntries(),
			readDirError:    nil,
			removeError:     nil,
			symlinkError:    nil,
			chmodError:      nil,
			expectedError:   nil,
		},
		{
			testTitle:       "should fail when reading the path returns an error",
			readDirResponse: getEmptyDirEntries(),
			readDirError:    errors.New("an error occurred while reading the give path"),
			removeError:     nil,
			symlinkError:    nil,
			chmodError:      nil,
			expectedError:   errors.New("an error occurred while reading the give path"),
		},
		{
			testTitle:       "should fail when deleting the existing symlink returns an error",
			readDirResponse: getGoVersionBinEntries(),
			readDirError:    nil,
			removeError:     errors.New("an error occurred while deleting the symlink"),
			symlinkError:    nil,
			chmodError:      nil,
			expectedError:   errors.New("an error occurred while deleting the symlink"),
		},
		{
			testTitle:       "should fail when creating the symlink returns an error",
			readDirResponse: getGoVersionBinEntries(),
			readDirError:    nil,
			removeError:     nil,
			symlinkError:    errors.New("an error occurred while creating the symlink"),
			chmodError:      nil,
			expectedError:   errors.New("an error occurred while creating the symlink"),
		},
		{
			testTitle:       "should fail when creating the mode change returns an error",
			readDirResponse: getGoVersionBinEntries(),
			readDirError:    nil,
			removeError:     nil,
			symlinkError:    nil,
			chmodError:      errors.New("an error occurred while updating the permissions"),
			expectedError:   errors.New("an error occurred while updating the permissions"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:             "/tmp",
				ReadDirMockResponse: getDirEntries(tc.readDirResponse),
				ReadDirError:        tc.readDirError,
				RemoveError:         tc.removeError,
				SymlinkError:        tc.symlinkError,
				ChmodError:          tc.chmodError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			err := fileHelper.CreateExecutableSymlink("go1.21.0")

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}
		})
	}
}

func TestUpdateRecentVersion(t *testing.T) {
	testCases := []struct {
		testTitle        string
		createError      error
		writeToFileError error
		expectedError    error
	}{
		{
			testTitle:        "should update the recent file content",
			createError:      nil,
			writeToFileError: nil,
			expectedError:    nil,
		},
		{
			testTitle:        "should return an error when the file creation fails",
			createError:      errors.New("an error occurred while creating the file"),
			writeToFileError: nil,
			expectedError:    errors.New("an error occurred while creating the file"),
		},
		{
			testTitle:        "should return an error when the file content update fails",
			createError:      nil,
			writeToFileError: errors.New("an error occurred while writing on the file"),
			expectedError:    errors.New("an error occurred while writing on the file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:          "/tmp",
				CreateError:      tc.createError,
				WriteStringError: tc.writeToFileError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			err := fileHelper.UpdateRecentVersion("go1.21.0")

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}
		})
	}
}

func TestStoreVersionsResponse(t *testing.T) {
	testCases := []struct {
		testTitle        string
		writeToFileError error
	}{
		{
			testTitle:        "should store the new version to the file",
			writeToFileError: nil,
		},
		{
			testTitle:        "should store the new version to the file",
			writeToFileError: errors.New("an error occurred while writing to the file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:        "/tmp",
				WriteFileError: tc.writeToFileError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			err := fileHelper.StoreVersionsResponse([]byte("Some versions JSON"))

			if tc.writeToFileError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
			}

			if tc.writeToFileError != nil && err.Error() != tc.writeToFileError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.writeToFileError.Error(), err.Error())
				return
			}
		})
	}
}

func TestGetCachedResponse(t *testing.T) {
	testCases := []struct {
		testTitle       string
		readFileError   error
		readFileContent any
		expectedError   error
	}{
		{
			testTitle:     "should return caches response JSON",
			readFileError: nil,
			readFileContent: []map[string]interface{}{
				{
					"version": "go1.21.0",
					"stable":  true,
					"files": []any{
						map[string]any{
							"arch":     string("arm64"),
							"filename": string("go1.21.0.linux-arm64.tar.gz"),
							"kind":     string("archive"),
							"os":       string("linux"),
							"sha256":   string("818d46ede85682dd551ad378ef37a4d247006f12ec59b5b755601d2ce114369a"),
							"size":     float64(9.6962473e+07),
							"version":  string("go1.21.0"),
						},
						map[string]any{
							"arch":     string("amd64"),
							"filename": string("go1.21.0.darwin-amd64.pkg"),
							"kind":     string("archive"),
							"os":       string("darwin"),
							"sha256":   string("725de310e4cba0121d6337053b2cfc3fe47da4a3d50726582731cb1d2a70137e"),
							"size":     float64(6.714125e+07),
							"version":  string("go1.21.0"),
						},
					},
				}},
			expectedError: nil,
		},
		{
			testTitle:       "should return an error when reading of the file fails",
			readFileError:   errors.New("an error occurred while readint the file"),
			readFileContent: nil,
			expectedError:   errors.New("an error occurred while readint the file"),
		},
		{
			testTitle:       "should return an error when unmarhsalling the response fails",
			readFileError:   nil,
			readFileContent: bytes.NewBufferString("{foo: bar}"), // force syntax error to response body
			expectedError:   errors.New("json: cannot unmarshal object into Go value of type []api_client.VersionInfo"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			versions := tc.readFileContent
			versionsBody, _ := json.Marshal(versions)

			fs := testutils.FakeFileSystem{
				HomeDir:       "/tmp",
				ReadFileBytes: versionsBody,
				ReadFileError: tc.readFileError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})

			var responseVersions []api_client.VersionInfo
			err := fileHelper.GetCachedResponse(&responseVersions)

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}

			if tc.expectedError == nil && len(responseVersions) == 0 {
				t.Errorf("response should not be an empty array")
				return
			}
		})
	}
}

func TestAreVersionsCached(t *testing.T) {
	testCases := []struct {
		testTitle      string
		startResponse  testutils.FakeFileInfo
		statError      error
		diffInHours    float64
		expectedResult bool
	}{
		{
			testTitle:      "should return false when the stat returns an error",
			startResponse:  testutils.FakeFileInfo{},
			statError:      errors.New("an error occurred when getting the stats for the given file"),
			diffInHours:    0,
			expectedResult: false,
		},
		{
			testTitle: "should return false when the file is older than a week",
			startResponse: testutils.FakeFileInfo{
				FileName:    ".go.versions",
				FileSize:    1000,
				FileMode:    0,
				FileModTime: time.Time{},
				FileIsDir:   false,
			},
			statError:      nil,
			diffInHours:    24 * 8,
			expectedResult: false,
		},
		{
			testTitle: "should return true when the file is not older than a week",
			startResponse: testutils.FakeFileInfo{
				FileName:    ".go.versions",
				FileSize:    1000,
				FileMode:    0,
				FileModTime: time.Time{},
				FileIsDir:   false,
			},
			statError:      nil,
			diffInHours:    24,
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:          "/tmp",
				StatError:        tc.statError,
				StatMockResponse: tc.startResponse,
			}
			fakeClock := testutils.FakeClock{GetDiffInHoursFromNowRes: tc.diffInHours}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, fakeClock)
			res := fileHelper.AreVersionsCached()

			if res != tc.expectedResult {
				t.Errorf("response should be %t, instead got %t", tc.expectedResult, res)
			}
		})
	}
}

func TestGetRecentVersion(t *testing.T) {
	testCases := []struct {
		testTitle      string
		readFileError  error
		fileContent    []byte
		expectedResult string
	}{
		{
			testTitle:      "should return the content from the file",
			readFileError:  nil,
			fileContent:    []byte("go1.20.7"),
			expectedResult: "go1.20.7",
		},
		{
			testTitle:      "should return an empty string when an error has occurred and log the error",
			readFileError:  errors.New("Some error occurred while reading the file"),
			fileContent:    []byte("go1.20.7"),
			expectedResult: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:       "/tmp",
				ReadFileBytes: tc.fileContent,
				ReadFileError: tc.readFileError,
			}

			logWriter := &testutils.FakeStdout{}
			fileHelper := createFileHelper(&testutils.FakeStdout{}, logWriter, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			res := fileHelper.GetRecentVersion()

			if res != tc.expectedResult {
				t.Errorf("result should be %q, instead got %q", tc.expectedResult, res)
				return
			}

			if tc.readFileError != nil {
				errorMessage := tc.readFileError.Error() + "\n"
				logMessage := logWriter.GetPrintMessages()[0]
				regex := fmt.Sprintf("(ERROR:) (\\d{4}\\/\\d{2}\\/\\d{2}) (\\d{2}:\\d{2}:([0-9]*[.])?[0-9]+) .+: %s", errorMessage)
				match, _ := regexp.MatchString(regex, logMessage)
				if !match {
					t.Errorf("Wrong log messages received, got=%s", logMessage)
				}
			} else {
				if len(logWriter.GetPrintMessages()) != 0 {
					t.Errorf("no logs should be printed, instead got %v", logWriter.GetPrintMessages())
				}
			}
		})
	}
}

func TestDirectoryExists(t *testing.T) {
	testCases := []struct {
		testTitle      string
		statError      error
		expectedResult bool
	}{
		{
			testTitle:      "should return true when the directory exists",
			statError:      nil,
			expectedResult: true,
		},
		{
			testTitle:      "should return false when start returns any error",
			statError:      errors.New("file does not exist"),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:   "/tmp",
				StatError: tc.statError,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			res := fileHelper.DirectoryExists("go1.21.0")

			if res != tc.expectedResult {
				t.Errorf("res should be %t, instead for %t", tc.expectedResult, res)
			}
		})
	}
}

func TestDeleteDirectory(t *testing.T) {
	testCases := []struct {
		testTitle      string
		removeAllError error
	}{
		{
			testTitle:      "should remove the directory",
			removeAllError: nil,
		},
		{
			testTitle:      "should return an error when deleting fails",
			removeAllError: errors.New("an error occurred while deleting the directory"),
		},
	}

	for _, tc := range testCases {
		fs := testutils.FakeFileSystem{
			HomeDir:        "/tmp",
			RemoveAllError: tc.removeAllError,
		}

		fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
		err := fileHelper.DeleteDirectory("go1.21.0")

		if tc.removeAllError == nil && err != nil {
			t.Errorf("error should be nil, instead got %q", err.Error())
			return
		}

		if tc.removeAllError != nil && err.Error() != tc.removeAllError.Error() {
			t.Errorf("error should be %q, instead got %q", tc.removeAllError.Error(), err.Error())
			return
		}
	}
}

func TestCreateInitFiles(t *testing.T) {
	testCases := []struct {
		testTitle     string
		pathToFail    string
		mkdirError    error
		openFileError error
		expectedError error
	}{
		{
			testTitle:     "should create the initial files",
			pathToFail:    "",
			mkdirError:    nil,
			openFileError: nil,
			expectedError: nil,
		},
		{
			testTitle:     "should fail when .gvs can't be created",
			pathToFail:    "/tmp/.gvs",
			mkdirError:    errors.New("error while creating the .gvs file"),
			openFileError: nil,
			expectedError: errors.New("error while creating the .gvs file"),
		},
		{
			testTitle:     "should fail when .gvs/.go.versions can't be created",
			pathToFail:    "/tmp/.gvs/.go.versions",
			mkdirError:    errors.New("error while creating the .go.versions file"),
			openFileError: nil,
			expectedError: errors.New("error while creating the .go.versions file"),
		},
		{
			testTitle:     "should fail when /tmp/bin can't be created",
			pathToFail:    "/tmp/bin",
			mkdirError:    errors.New("error while creating the /tmp/bin file"),
			openFileError: nil,
			expectedError: errors.New("error while creating the /tmp/bin file"),
		},
		{
			testTitle:     "should fail when /tmp/.gvs/gvs.log can't be created",
			pathToFail:    "",
			mkdirError:    nil,
			openFileError: errors.New("error while creating the /tmp/.gvs/gvs.log file"),
			expectedError: errors.New("error while creating the /tmp/.gvs/gvs.log file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			var logFile *os.File

			fs := testutils.FakeFileSystem{
				HomeDir:                   "/tmp",
				MkdirIfNotExistPathToFail: tc.pathToFail,
				MkdirIfNotExistError:      tc.mkdirError,
				OpenFileError:             tc.openFileError,
				OpenFileMockFile:          logFile,
			}

			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, testutils.FakeClock{})
			returnedFile, err := fileHelper.CreateInitFiles()

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}

			if tc.expectedError == nil && logFile != returnedFile {
				t.Error("the created file and the returned file addresses do not match")
				return
			}

			if tc.expectedError != nil && returnedFile != nil {
				t.Error("the returned file should be nil")
				return
			}
		})
	}
}

func TestGetLatestCreatedGoVersionDirectory(t *testing.T) {
	testCases := []struct {
		testTitle       string
		readDirError    error
		readDirResponse []testutils.FakeDirEntry
		expectedError   error
		expectedResult  string
	}{
		{
			testTitle:       "should fail when ReadDir returns an error back",
			readDirError:    errors.New("an error occurred while reading the directory path"),
			readDirResponse: getEmptyDirEntries(),
			expectedError:   errors.New("an error occurred while reading the directory path"),
			expectedResult:  "",
		},
		{
			testTitle:       "should return an empty string when no file entry is a directory",
			readDirError:    nil,
			readDirResponse: getGoVersionBinEntries(),
			expectedError:   nil,
			expectedResult:  "",
		},
		{
			testTitle:       "should failt when file info returns an error",
			readDirError:    nil,
			readDirResponse: getGoDirEntry(errors.New("some error while getting ")),
			expectedError:   errors.New("some error while getting "),
			expectedResult:  "",
		},
		{
			testTitle: "should return the latest created directory name", // also test the case when the latest is a file
			readDirResponse: []testutils.FakeDirEntry{
				{
					DirEntryName:  "go",
					DirEntryIsDir: true,
					DireEntryType: 0,
					DirEntryInfo: testutils.FakeFileInfo{
						FileName:    "go",
						FileSize:    1000,
						FileMode:    0,
						FileModTime: time.Date(2021, 10, 21, 10, 30, 0, 0, time.Local),
						FileIsDir:   true,
					},
					DirEntryInfoError: nil,
				},
				{
					DirEntryName:  "some-file",
					DirEntryIsDir: false,
					DireEntryType: 0,
					DirEntryInfo: testutils.FakeFileInfo{
						FileName:    "some-file",
						FileSize:    1000,
						FileMode:    0,
						FileModTime: time.Date(2021, 10, 21, 11, 30, 0, 0, time.Local),
						FileIsDir:   false,
					},
					DirEntryInfoError: nil,
				},
				{
					DirEntryName:  "go1.20.10",
					DirEntryIsDir: true,
					DireEntryType: 0,
					DirEntryInfo: testutils.FakeFileInfo{
						FileName:    "go1.20.10",
						FileSize:    1000,
						FileMode:    0,
						FileModTime: time.Date(2021, 10, 21, 7, 30, 0, 0, time.Local),
						FileIsDir:   true,
					},
					DirEntryInfoError: nil,
				},
			},
			readDirError:   nil,
			expectedError:  nil,
			expectedResult: "go",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			fs := testutils.FakeFileSystem{
				HomeDir:             "/tmp",
				ReadDirError:        tc.readDirError,
				ReadDirMockResponse: getDirEntries(tc.readDirResponse),
			}

			clock := testutils.FakeClock{
				UseRealIsBefore: true,
				UseRealIsAfter:  true,
			}
			fileHelper := createFileHelper(&testutils.FakeStdout{}, nil, fs, testutils.FakeUnzipper{}, clock)
			dirName, err := fileHelper.GetLatestCreatedGoVersionDirectory()

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}

			if dirName != tc.expectedResult {
				t.Errorf("the directory name should be %q, instead got %q", tc.expectedResult, dirName)
				return
			}
		})
	}
}
