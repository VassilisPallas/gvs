package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GoRepositoryAPI interface {
	FetchVersions(url string, v *[]VersionInfo) error
	DownloadVersion(url string, cb func(body io.Reader) error) error
}

type Go struct {
	client HTTPClient
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

func (g *Go) FetchVersions(url string, v *[]VersionInfo) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/?mode=json&include=all", url), nil)
	if err != nil {
		return err
	}

	response, err := g.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status %d", response.StatusCode)
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

func (g *Go) DownloadVersion(url string, cb func(body io.Reader) error) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := g.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Request failed with status %d", response.StatusCode)
	}

	if err := cb(response.Body); err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func New() GoRepositoryAPI {
	return &Go{client: &http.Client{}}
}
