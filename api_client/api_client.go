package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VassilisPallas/gvs/config"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GoClientAPI interface {
	FetchVersions(ctx context.Context, v *[]VersionInfo) error
	DownloadVersion(ctx context.Context, filename string, cb func(body io.ReadCloser) error) error
}

type Go struct {
	Client HTTPClient
	Config config.Configuration

	GoClientAPI
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

// TODO: add more tests for config, like url, timeout and ctx cancel etc

func (g Go) FetchVersions(ctx context.Context, v *[]VersionInfo) error {
	// TODO: write test for that
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/?mode=json&include=all", g.Config.GO_BASE_URL), nil)
	if err != nil {
		return err
	}

	response, err := g.Client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", response.StatusCode)
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
	// TODO: write test for that
	url := fmt.Sprintf("%s/%s", g.Config.GO_BASE_URL, filename)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := g.Client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("request failed with status %d", response.StatusCode)
	}

	if err := cb(response.Body); err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// TODO: check if should return *Go
func New(config config.Configuration) Go {
	return Go{Client: &http.Client{
		Timeout: time.Duration(config.REQUEST_TIMEOUT) * time.Second,
	}, Config: config}
}
