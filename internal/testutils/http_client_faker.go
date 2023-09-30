package testutils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/VassilisPallas/gvs/api_client"
)

type MockClient struct {
	Status       int
	Body         io.ReadCloser
	RequestError error
}

func (tc MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: tc.Status,
		Body:       tc.Body,
	}, tc.RequestError
}

type FakeGoClientAPI struct {
	DownloadError      error
	FetchVersionsError error
}

func (ga FakeGoClientAPI) FetchVersions(ctx context.Context, v *[]api_client.VersionInfo) error {
	if ga.FetchVersionsError != nil {
		return ga.FetchVersionsError
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
	}

	responseBody, err := json.Marshal(responseVersions)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &v)

	return err
}

func (ga FakeGoClientAPI) DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error {
	if err := cb(nil); err != nil {
		return err
	}

	return ga.DownloadError
}
