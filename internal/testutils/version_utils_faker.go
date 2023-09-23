package testutils

import (
	"encoding/json"
	"fmt"

	"github.com/VassilisPallas/gvs/api_client"
)

type FakeVersionHelper struct {
	CachedVersion        bool
	CacheResponseError   error
	RecentVersion        string
	DeleteDirectoryError error
}

func (FakeVersionHelper) StoreVersionsResponse(body []byte) error {
	return nil
}

func (vh FakeVersionHelper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	if vh.CacheResponseError != nil {
		return vh.CacheResponseError
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
	fmt.Println(err)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &v)

	return err
}

func (vh FakeVersionHelper) AreVersionsCached() bool {
	return vh.CachedVersion
}

func (vh FakeVersionHelper) GetRecentVersion() string {
	return vh.RecentVersion
}

func (FakeVersionHelper) DirectoryExists(goVersion string) bool {
	return false
}

func (vh FakeVersionHelper) DeleteDirectory(dirName string) error {
	if dirName == "bad_version" {
		return vh.DeleteDirectoryError
	}

	return nil
}
