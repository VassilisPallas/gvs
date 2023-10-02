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

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBuffer(responseBody)},
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err != nil {
		t.Errorf("FetchVersions error should be nil, instead got %q", err)
	}

	if versions == nil {
		t.Error("FetchVersions versions shouldn't be nil")
	}

	var responseToMap []map[string]interface{}
	inrec, _ := json.Marshal(&versions)
	json.Unmarshal(inrec, &responseToMap)

	if !cmp.Equal(responseToMap, responseVersions) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(responseVersions, responseToMap))
	}
}

func TestFetchVersionNewRequestWithContextError(t *testing.T) {
	expectedError := fmt.Errorf("parse \" https://go.dev/dl/?mode=json&include=all\": first path segment in URL cannot contain colon")

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("")},
	}

	goRepo := api_client.New(client, " https://go.dev/dl") // add space to URL to raise an error

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.TODO(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %q, instead got %q", expectedError.Error(), err.Error())
	}

	if versions != nil {
		t.Error("FetchVersions versions should be nil")
		return
	}
}

func TestFetchVersionsNonOkStatus(t *testing.T) {
	expectedError := fmt.Errorf("request failed with status %d", http.StatusBadRequest)

	client := testutils.MockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestFetchVersionsRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")

	client := testutils.MockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestFetchVersionsUnmarshalFailed(t *testing.T) {
	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("{foo: bar}")}, // force syntax error to response body
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	expectedError := fmt.Errorf("invalid character 'f' looking for beginning of object key string")
	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestFetchVersionsReadBodyFailed(t *testing.T) {
	expectedError := fmt.Errorf("some error while reading body")

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopReaderCloser{readError: expectedError, Reader: bytes.NewBufferString("")},
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	var versions []api_client.VersionInfo
	err := goRepo.FetchVersions(context.Background(), &versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("FetchVersions error should be %q, instead got %q", expectedError.Error(), err.Error())
		return
	}
}

func TestDownloadVersionSuccess(t *testing.T) {
	cb := func(r io.ReadCloser) error {
		return nil
	}

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err != nil {
		t.Errorf("FetchVersions error should be nil, instead got %q", err.Error())
		return
	}
}

func TestDownloadVersionNewRequestWithContextError(t *testing.T) {
	expectedError := fmt.Errorf("parse \" https://go.dev/dl/some_file_name\": first path segment in URL cannot contain colon")

	cb := func(r io.ReadCloser) error {
		return nil
	}

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}

	goRepo := api_client.New(client, " https://go.dev/dl") // add space to URL to raise an error

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("DownloadVersion error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestDownloadVersionNonOkStatus(t *testing.T) {
	expectedError := fmt.Errorf("request failed with status %d", http.StatusBadRequest)
	cb := func(r io.ReadCloser) error {
		return nil
	}

	client := testutils.MockClient{
		Status: http.StatusBadRequest,
		Body:   nil,
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("DownloadVersion error should be %q, instead got %q", expectedError.Error(), err.Error())
		return
	}
}

func TestDownloadVersionRequestFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error")
	cb := func(r io.ReadCloser) error {
		return nil
	}

	client := testutils.MockClient{
		Status:       http.StatusBadRequest,
		Body:         nil,
		RequestError: expectedError,
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("DownloadVersion error should be %q, instead got %q", expectedError.Error(), err.Error())
		return
	}
}

func TestDownloadVersionCallbackFailed(t *testing.T) {
	expectedError := fmt.Errorf("Some error within the callback")
	cb := func(r io.ReadCloser) error {
		return expectedError
	}

	client := testutils.MockClient{
		Status: http.StatusOK,
		Body:   nopCloser{bytes.NewBufferString("foo")},
	}

	goRepo := api_client.New(client, "https://go.dev/dl")

	err := goRepo.DownloadVersion(context.Background(), "some_file_name", cb)

	if err.Error() != expectedError.Error() {
		t.Errorf("DownloadVersion error should be %q, instead got %q", expectedError.Error(), err.Error())
		return
	}
}
