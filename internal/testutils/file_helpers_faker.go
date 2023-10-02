package testutils

import (
	"encoding/json"
	"io"
	"os"

	"slices"

	"github.com/VassilisPallas/gvs/api_client"
)

type FakeFilesHelper struct {
	CreateExecutableSymlinkError error
	UpdateRecentVersionError     error
	CreateTarFileError           error
	UnzippingError               error
	RenameDirectoryError         error
	RemoveTarFileError           error
	CacheResponseError           error
	DeleteDirectoryError         error
	GetTarChecksumError          error

	Checksum                  string
	RecentVersion             string
	CachedVersion             bool
	AlreadyDownloadedVersions []string

	RemoveTarFileCalled bool
}

func (fh FakeFilesHelper) CreateTarFile(content io.ReadCloser) error {
	return fh.CreateTarFileError
}

func (fh *FakeFilesHelper) GetTarChecksum() (string, error) {
	return fh.Checksum, fh.GetTarChecksumError
}

func (fh FakeFilesHelper) UnzipTarFile() error {
	return fh.UnzippingError
}

func (fh FakeFilesHelper) RenameGoDirectory(goVersionName string) error {
	return fh.RenameDirectoryError
}

func (fh *FakeFilesHelper) RemoveTarFile() error {
	fh.RemoveTarFileCalled = true
	return fh.RemoveTarFileError
}

func (fh FakeFilesHelper) CreateExecutableSymlink(goVersionName string) error {
	return fh.CreateExecutableSymlinkError
}

func (fh FakeFilesHelper) UpdateRecentVersion(goVersionName string) error {
	return fh.UpdateRecentVersionError
}

func (FakeFilesHelper) StoreVersionsResponse(body []byte) error {
	return nil
}

func (fh FakeFilesHelper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	if fh.CacheResponseError != nil {
		return fh.CacheResponseError
	}

	responseVersions := []map[string]interface{}{
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
		},
		{
			"version": "go1.20.0",
			"stable":  true,
			"files":   []any{},
		},
		{
			"version": "go1.19.0",
			"stable":  true,
			"files":   []any{},
		},
		{
			"version": "go1.18.0",
			"stable":  true,
			"files":   []any{},
		},
	}

	responseBody, err := json.Marshal(responseVersions)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &v)

	return err
}

func (fh FakeFilesHelper) AreVersionsCached() bool {
	return fh.CachedVersion
}

func (fh FakeFilesHelper) GetRecentVersion() string {
	return fh.RecentVersion
}

func (fh FakeFilesHelper) DirectoryExists(goVersion string) bool {
	return slices.Contains(fh.AlreadyDownloadedVersions, goVersion)
}

func (fh FakeFilesHelper) DeleteDirectory(dirName string) error {
	if dirName == "bad_version" {
		return fh.DeleteDirectoryError
	}

	return nil
}

func (fh FakeFilesHelper) CreateInitFiles() (*os.File, error) {
	return nil, nil
}
