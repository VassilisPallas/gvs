// Package api_client provides an interface to make request for
// fetching and downloading versions.
package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VassilisPallas/gvs/errors"
)

// HTTPClient is an interface to let the user of the package
// to use any client they want to use to make the requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GoClientAPI is the interface that wraps the basic methods for making requests.
//
// FetchVersions fetches and returns the available Go versions.
// The versions should be parsed and stored in the value pointed to by v.
// FetchVersions must return a non-null error if the request, or the parsing of the response fails.
//
// DownloadVersion downloads the content (most likely a tar.gz file) and then is passing
// the response to the callack function.
// DownloadVersion must close the response body reader after passing it in the callback function.
// DownloadVersion must return an non-null error if the request failes or the callback returns
// an non-null error.
type GoClientAPI interface {
	FetchVersions(ctx context.Context, v *[]VersionInfo) error
	DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error
}

// Go is the struct that implements the GoClientAPI interface
type Go struct {
	// Client will be used as a custom HTTPClient to make the request and return the response.
	client HTTPClient
	// baseURL will be used from both FetchVersions and  DownloadVersion methods.
	baseURL string
}

// VersionInfo is the struct that represents the response JSON for the versions.
type VersionInfo struct {
	// Version contains the Go version (e.g. `go1.21.3`)
	Version string `json:"version"`
	// IsStable is boolean if the Go version is stable or not. Unstable versions are the release candidates (e.g. `go1.21rc4`).
	IsStable bool `json:"stable"`
	// Files is a slice that contains information for different kind of files that are available for the given version,
	// depending on the OS and Architecture type.
	Files []FileInformation `json:"files"`
}

// FileInformation is the struct that represents the JSON for the files inside the versions.
type FileInformation struct {
	// Filename contains the name of the archived file. This is used as a parameter on the DownloadVersion method.
	Filename string `json:"filename"`
	// OS contains the type of the Operating System (e.g. `darwin`, `linux`, `windows` etc).
	OS string `json:"os"`
	// Architecture contains the architecture target (e.g. `386`, `amd64`, `arm64`, `s390x` etc).
	Architecture string `json:"arch"`
	// Version contains the Go version (e.g. `go1.21.3`).
	Version string `json:"version"`
	// Checksum contains the SHA256 Checksum for the given file.
	Checksum string `json:"sha256"`
	// Size contains the given file in bytes.
	Size uint64 `json:"size"`
	// The Kind represents the kind of the file:
	// one of source, archive, or installer.
	Kind string `json:"kind"`
}

// FetchVersions fetches and returns the available Go versions.
//
// It is using the NewRequestWithContext function from the `http` package.
// If the request is successful it will parses the JSON-encoded data and store it
// in the value pointed to by v.
//
// The ctx parameter is used for the request. It can be any type of context
// depending on the use case.
//
// If the request or the parse of the response body fails,
// FetchVersions will return an error.
func (g Go) FetchVersions(ctx context.Context, v *[]VersionInfo) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/?mode=json&include=all", g.baseURL), nil)
	if err != nil {
		return err
	}

	response, err := g.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return &errors.RequestError{StatusCode: response.StatusCode}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}

	return nil
}

// DownloadVersion downloads the content (most likely a tar.gz file) and then is passing
// the response to the callack function.
//
// It is using the NewRequestWithContext function from the `http` package.
// If the request is successful it will pass the response body to the callback.
//
// The ctx parameter is used for the request. It can be any type of context
// depending on the use case.
//
// If the request failes or the callback returns an non-null error,
// DownloadVersion will return an error.
//
// DownloadVersion finally closed response body reader after the execution of the method.
func (g Go) DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error {
	url := fmt.Sprintf("%s/%s", g.baseURL, filename)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := g.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return &errors.RequestError{StatusCode: response.StatusCode}
	}

	if err := cb(response.Body); err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// New returns a Go instance that implements the GoClientAPI interface.
// Each call to New returns a distinct Go instance even if the parameters are identical.
func New(client HTTPClient, baseURL string) Go {
	return Go{client: client, baseURL: baseURL}
}
