package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockClient struct {
	Status       int
	Body         io.ReadCloser
	RequestError error
}

func (tc mockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: tc.Status,
		Body:       tc.Body,
	}, tc.RequestError
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestFetchVersionsSuccess(t *testing.T) {
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
		}}

	responseBody, _ := json.Marshal(responseVersions)

	goRepo := &Go{client: mockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBuffer(responseBody)},
	}}

	var versions []VersionInfo
	err := goRepo.FetchVersions("someurl", &versions)

	if err != nil {
		t.Errorf("FetchVersions error should be nil, instead got %s", err)
		return
	}

	if versions == nil {
		t.Error("FetchVersions versions shouldn't be nil")
		return
	}

	var responseToMap []map[string]interface{}
	inrec, _ := json.Marshal(&versions)
	json.Unmarshal(inrec, &responseToMap)

	if !cmp.Equal(responseToMap, responseVersions) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(responseVersions, responseToMap))
	}
}

func TestFetchVersionsNonOkStatus(t *testing.T) {
	expectedError := fmt.Errorf("Request failed with status %d", http.StatusBadRequest)

	goRepo := &Go{client: mockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}}

	var versions []VersionInfo
	err := goRepo.FetchVersions("someurl", &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestFetchVersionsRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")

	goRepo := &Go{client: mockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}}

	var versions []VersionInfo
	err := goRepo.FetchVersions("someurl", &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestDownloadVersionSuccess(t *testing.T) {
	cb := func(r io.Reader) error {
		return nil
	}

	goRepo := &Go{client: mockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}}

	err := goRepo.DownloadVersion("someurl", cb)

	if err != nil {
		t.Errorf("FetchVersions error should be nil, instead got %s", err)
		return
	}
}

func TestDownloadVersionNonOkStatus(t *testing.T) {
	expectedError := fmt.Errorf("Request failed with status %d", http.StatusBadRequest)
	cb := func(r io.Reader) error {
		return nil
	}

	goRepo := &Go{client: mockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}}

	err := goRepo.DownloadVersion("someurl", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestDownloadVersionRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")
	cb := func(r io.Reader) error {
		return nil
	}

	goRepo := &Go{client: mockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}}

	err := goRepo.DownloadVersion("someurl", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestDownloadVersionCallbackFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error within the callback")
	cb := func(r io.Reader) error {
		return expectedError
	}

	goRepo := &Go{client: mockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}}

	err := goRepo.DownloadVersion("someurl", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}
