package api_client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/google/go-cmp/cmp"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type nopReaderCloser struct {
	readError error

	io.Reader
}

func (nrc nopReaderCloser) Read(p []byte) (n int, err error) {
	return 0, nrc.readError
}

func (nopReaderCloser) Close() error { return nil }

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

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBuffer(responseBody)},
	}}

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

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

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}}

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestFetchVersionsRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}}

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestFetchVersionsUnmarshalFailed(t *testing.T) {
	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("{foo: bar}")}, // force syntax error to response body
	}}

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err == nil {
		t.Errorf("FetchVersions error should be '%s', instead got nil", err)
		return
	}
}

func TestFetchVersionsReadBodyFailed(t *testing.T) {
	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopReaderCloser{readError: fmt.Errorf("some error while reading body"), Reader: bytes.NewBufferString("")},
	}}

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err == nil {
		t.Errorf("FetchVersions error should be '%s', instead got nil", err)
		return
	}
}

func TestDownloadVersionSuccess(t *testing.T) {
	cb := func(r io.ReadCloser) error {
		return nil
	}

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}}

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err != nil {
		t.Errorf("FetchVersions error should be nil, instead got %s", err)
		return
	}
}

func TestDownloadVersionNonOkStatus(t *testing.T) {
	expectedError := fmt.Errorf("Request failed with status %d", http.StatusBadRequest)
	cb := func(r io.ReadCloser) error {
		return nil
	}

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}}

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestDownloadVersionRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")
	cb := func(r io.ReadCloser) error {
		return nil
	}

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}}

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}

func TestDownloadVersionCallbackFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error within the callback")
	cb := func(r io.ReadCloser) error {
		return expectedError
	}

	goRepo := &api_client.Go{Client: testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}}

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %s, instead got %s", expectedError, err)
		return
	}
}
