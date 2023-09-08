package version_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/version"
)

func TestStoreVersionsResponse(t *testing.T) {
	appDir := "/tmp"
	versionResponseFile := "goVersions.json"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	err := helper.StoreVersionsResponse([]byte("Some content"))

	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}
}

func TestStoreVersionsResponseShouldNotFailWithEmptyContent(t *testing.T) {
	appDir := "/tmp"
	versionResponseFile := "goVersions.json"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	err := helper.StoreVersionsResponse([]byte(""))

	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}
}

func TestStoreVersionsResponseFileToPathThatDoesNotExist(t *testing.T) {
	appDir := "/some_other_dst"
	versionResponseFile := "goVersions.json"
	expectedError := "open /some_other_dst/goVersions.json: no such file or directory"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	err := helper.StoreVersionsResponse([]byte("Some content"))

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}
}

func TestGetCachedResponse(t *testing.T) {
	appDir := "/tmp"
	versionResponseFile := "goVersions.json"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	responseVersionsMap := []map[string]interface{}{
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
		}}

	responseBody, _ := json.Marshal(responseVersionsMap)

	helper.StoreVersionsResponse([]byte(responseBody))

	var responseVersions []api_client.VersionInfo
	err := helper.GetCachedResponse(&responseVersions)

	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}

	if responseVersions == nil {
		t.Errorf("responseVersions should not be nil")
	}
}

func TestGetCachedResponseFileToPathThatDoesNotExist(t *testing.T) {
	appDir := "/some_other_dst"
	versionResponseFile := "goVersions.json"
	expectedError := "open /some_other_dst/goVersions.json: no such file or directory"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	var responseVersions []api_client.VersionInfo
	err := helper.GetCachedResponse(&responseVersions)

	if err.Error() != expectedError {
		t.Errorf("error should be '%s', instead got '%s'", expectedError, err.Error())
	}
}

func TestGetCachedResponseUnmarshalFailed(t *testing.T) {
	appDir := "/tmp"
	versionResponseFile := "goVersions.json"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			AppDir:                  appDir,
			VersionResponseFileName: versionResponseFile,
		},
	}

	helper.StoreVersionsResponse([]byte("{foo: bar}")) // force syntax error to response body

	var responseVersions []api_client.VersionInfo
	err := helper.GetCachedResponse(&responseVersions)

	defer os.Remove(fmt.Sprintf("%s/%s", appDir, versionResponseFile))

	if err == nil {
		t.Errorf("error should be '%s', instead got nil", err)
	}
}

func TestGetRecentVersion(t *testing.T) {
	fileLocation := "/tmp/CURRENT"
	goVersion := "go1.21.0"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			CurrentVersionFile: fileLocation,
		},
	}

	file, _ := os.Create(fileLocation)
	io.WriteString(file, goVersion)

	content := helper.GetRecentVersion()

	if content != goVersion {
		t.Errorf("content is not correct. Got '%s', instead of '%s'", content, goVersion)
	}
}

func TestGetRecentVersionFileToPathThatDoesNotExist(t *testing.T) {
	fileLocation := "/some_other_dst/CURRENT"

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			CurrentVersionFile: fileLocation,
		},
	}

	content := helper.GetRecentVersion()

	if content != "" {
		t.Errorf("content should be empty string, instead it got '%s'", content)
	}
}

func TestDirectoryExists(t *testing.T) {
	dirLocation := "/tmp"
	goVersion := "go1.21.0"

	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			VersionDir: dirLocation,
		},
	}

	os.MkdirAll(path, 0755)

	exists := helper.DirectoryExists(goVersion)

	defer os.RemoveAll(path)

	if !exists {
		t.Errorf("%s should exist", path)
	}
}

func TestDirectoryExistsDirectoryNotFound(t *testing.T) {
	dirLocation := "/tmp"
	goVersion := "go1.21.0"

	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			VersionDir: dirLocation,
		},
	}

	exists := helper.DirectoryExists(goVersion)

	if exists {
		t.Errorf("%s should not exist", path)
	}
}

func TestDeleteDirectory(t *testing.T) {
	dirLocation := "/tmp"
	goVersion := "go1.21.0"

	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			VersionDir: dirLocation,
		},
	}

	os.MkdirAll(path, 0755)

	err := helper.DeleteDirectory(goVersion)

	defer os.RemoveAll(path)

	if err != nil {
		t.Errorf("error should be nil, instead got '%s'", err.Error())
	}
}

func TestDeleteDirectoryDeleteFails(t *testing.T) {
	dirLocation := "/tmp"
	goVersion := "." // This is forcing os.RemoveAll to fail
	expectedErrorMessage := "RemoveAll /tmp/.: invalid argument"

	path := fmt.Sprintf("%s/%s", dirLocation, goVersion)

	helper := version.Helper{
		FileUtils: testutils.FakeFiler{
			VersionDir: dirLocation,
		},
	}

	os.MkdirAll(path, 0755)

	err := helper.DeleteDirectory(goVersion)

	defer os.RemoveAll(path)

	if err == nil {
		t.Errorf("error should be '%s, instead got nil", expectedErrorMessage)
	}
}
