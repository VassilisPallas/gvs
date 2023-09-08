package testutils

import (
	"context"
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
	DownloadError error
}

func (FakeGoClientAPI) FetchVersions(ctx context.Context, v *[]api_client.VersionInfo) error {
	return nil
}

func (ga FakeGoClientAPI) DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error {
	if err := cb(nil); err != nil {
		return err
	}

	return ga.DownloadError
}
