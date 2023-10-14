// Package files provides interfaces for reading
// and writing files.
package files

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
)

type FileHelpers interface {
	CreateTarFile(content io.ReadCloser) error
	GetTarChecksum() (string, error)
	UnzipTarFile() error
	RenameGoDirectory(goVersionName string) error
	RemoveTarFile() error
	CreateExecutableSymlink(goVersionName string) error
	UpdateRecentVersion(goVersionName string) error
	StoreVersionsResponse(body []byte) error
	GetCachedResponse(v *[]api_client.VersionInfo) error
	AreVersionsCached() bool
	GetRecentVersion() string
	DirectoryExists(goVersion string) bool
	DeleteDirectory(dirName string) error
	CreateInitFiles() (*os.File, error)
}

type Helper struct {
	fileUtils  FileUtils
	unzip      Unziper
	fileSystem FS
}

func (h Helper) CreateTarFile(content io.ReadCloser) error {
	file, err := h.fileSystem.Create(h.fileUtils.GetTarFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = h.fileSystem.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func (h Helper) GetTarChecksum() (string, error) {
	hasher := sha256.New()
	path := h.fileUtils.GetTarFile()

	f, err := h.fileSystem.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := h.fileSystem.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// TODO: add tests
func (h Helper) UnzipTarFile() error {
	source := h.fileUtils.GetTarFile()
	destination := h.fileUtils.GetVersionsDir()

	return h.unzip.UnzipSource(destination, source)
}

func (h Helper) RenameGoDirectory(goVersionName string) error {
	target := fmt.Sprintf("%s/%s", h.fileUtils.GetVersionsDir(), "go")

	if err := h.fileSystem.Rename(target, fmt.Sprintf("%s/%s", h.fileUtils.GetVersionsDir(), goVersionName)); err != nil {
		return err
	}

	return nil
}

func (h Helper) RemoveTarFile() error {
	err := h.fileSystem.Remove(h.fileUtils.GetTarFile())
	if err != nil {
		return err
	}
	return nil
}

// TODO: add tests
func (h Helper) CreateExecutableSymlink(goVersionName string) error {
	destination := fmt.Sprintf("%s/%s/bin", h.fileUtils.GetVersionsDir(), goVersionName)

	files, err := h.fileSystem.ReadDir(destination)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		newFile := fmt.Sprintf("%s/%s", destination, f.Name())
		link := fmt.Sprintf("%s/%s", h.fileUtils.GetBinDir(), f.Name())

		if _, err := h.fileSystem.Lstat(link); err == nil {
			// TODO: send to logger instead
			err := h.fileSystem.Remove(link)
			if err != nil {
				return err
			}
		}

		if err := h.fileSystem.Symlink(newFile, link); err != nil {
			return err
		}

		if err := h.fileSystem.Chmod(link, 0700); err != nil {
			return err
		}
	}

	return nil
}

func (h Helper) UpdateRecentVersion(goVersionName string) error {
	path := h.fileUtils.GetCurrentVersionFile()

	file, err := h.fileSystem.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, goVersionName)
	return err
}

func (h Helper) StoreVersionsResponse(body []byte) error {
	return h.fileSystem.WriteFile(fmt.Sprintf("%s/%s", h.fileUtils.GetAppDir(), h.fileUtils.GetVersionResponseFile()), body, 0644)
}

func (h Helper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	byte_versions, err := h.fileSystem.ReadFile(fmt.Sprintf("%s/%s", h.fileUtils.GetAppDir(), h.fileUtils.GetVersionResponseFile()))
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
	file := fmt.Sprintf("%s/%s", h.fileUtils.GetAppDir(), h.fileUtils.GetVersionResponseFile())

	if info, err := h.fileSystem.Stat(file); err == nil {
		currentTime := time.Now()
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week
		return currentTime.Sub(info.ModTime()).Hours() < 24*7
	}

	return false
}

func (h Helper) GetRecentVersion() string {
	path := h.fileUtils.GetCurrentVersionFile()

	content, err := h.fileSystem.ReadFile(path)
	if err != nil {
		// TODO: pass error as log to file
		return ""
	}
	return string(content)
}

func (h Helper) DirectoryExists(goVersion string) bool {
	target := h.fileUtils.GetVersionsDir()

	if _, err := h.fileSystem.Stat(fmt.Sprintf("%s/%s", target, goVersion)); !os.IsNotExist(err) {
		return true
	}

	return false
}

func (h Helper) DeleteDirectory(dirName string) error {
	target := h.fileUtils.GetVersionsDir()
	err := h.fileSystem.RemoveAll(fmt.Sprintf("%s/%s", target, dirName))
	if err != nil {
		return err
	}

	return nil
}

// TODO: add tests
func (h Helper) CreateInitFiles() (*os.File, error) {
	if err := h.fileSystem.MkdirIfNotExist(h.fileUtils.GetAppDir(), 0755); err != nil {
		return nil, err
	}
	if err := h.fileSystem.MkdirIfNotExist(h.fileUtils.GetVersionsDir(), 0755); err != nil {
		return nil, err
	}
	if err := h.fileSystem.MkdirIfNotExist(h.fileUtils.GetBinDir(), 0755); err != nil {
		return nil, err
	}

	// create log file
	filename := fmt.Sprintf("%s/%s", h.fileUtils.GetAppDir(), h.fileUtils.GetLogFile())
	return h.fileSystem.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}

func New(fileUtils FileUtils) *Helper {
	fs := FileSystem{}

	return &Helper{
		fileUtils:  fileUtils,
		unzip:      Unzip{fs: fs},
		fileSystem: fs,
	}
}
