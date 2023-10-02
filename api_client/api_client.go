package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VassilisPallas/gvs/errors"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GoClientAPI interface {
	FetchVersions(ctx context.Context, v *[]VersionInfo) error
	DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error
}

type Go struct {
	client  HTTPClient
	baseURL string
}

type FileInformation struct {
	Filename     string `json:"filename"`
	OS           string `json:"os"`
	Architecture string `json:"arch"`
	Version      string `json:"version"`
	Checksum     string `json:"sha256"`
	Size         uint64 `json:"size"`
	Kind         string `json:"kind"`
}

type VersionInfo struct {
	Version  string            `json:"version"`
	IsStable bool              `json:"stable"`
	Files    []FileInformation `json:"files"`
}

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

func New(client HTTPClient, baseURL string) Go {
	return Go{client: client, baseURL: baseURL}
}
