package version

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/files"
)

type VersionHelper interface {
	StoreVersionsResponse(body []byte) error
	GetCachedResponse(v *[]api_client.VersionInfo) error
	AreVersionsCached() bool
	GetRecentVersion() string
	DirectoryExists(goVersion string) bool
	DeleteDirectory(dirName string) error
}

type Helper struct {
	FileUtils files.FileUtils

	VersionHelper
}

func (h Helper) StoreVersionsResponse(body []byte) error {
	return os.WriteFile(fmt.Sprintf("%s/%s", h.FileUtils.GetAppDir(), h.FileUtils.GetVersionResponseFile()), body, 0644)
}

func (h Helper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	byte_versions, err := os.ReadFile(fmt.Sprintf("%s/%s", h.FileUtils.GetAppDir(), h.FileUtils.GetVersionResponseFile()))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(byte_versions, &v); err != nil {
		return err
	}

	return nil
}

// TODO: add tests
func (h Helper) AreVersionsCached() bool {
	file := fmt.Sprintf("%s/%s", h.FileUtils.GetAppDir(), h.FileUtils.GetVersionResponseFile())

	if info, err := os.Stat(file); err == nil {
		currentTime := time.Now()
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week
		return currentTime.Sub(info.ModTime()).Hours() > 24*7
	}

	return true
}

func (h Helper) GetRecentVersion() string {
	path := h.FileUtils.GetCurrentVersionFile()

	content, err := os.ReadFile(path)
	if err != nil {
		// TODO: pass error as log to file
		return ""
	}
	return string(content)
}

func (h Helper) DirectoryExists(goVersion string) bool {
	target := h.FileUtils.GetVersionsDir()

	if _, err := os.Stat(fmt.Sprintf("%s/%s", target, goVersion)); !os.IsNotExist(err) {
		return true
	}

	return false
}

func (h Helper) DeleteDirectory(dirName string) error {
	target := h.FileUtils.GetVersionsDir()
	err := os.RemoveAll(fmt.Sprintf("%s/%s", target, dirName))
	if err != nil {
		return err
	}

	return nil
}
