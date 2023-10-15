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
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/pkg/unzip"
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
	GetLatestCreatedGoVersionDirectory() (string, error)
}

type Helper struct {
	fileSystem FS
	unzip      unzip.Unzipper
	log        *logger.Log
}

func (h Helper) CreateTarFile(content io.ReadCloser) error {
	file, err := h.fileSystem.Create(getTarFile(h.fileSystem))
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
	path := getTarFile(h.fileSystem)

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

func (h Helper) UnzipTarFile() error {
	source := getTarFile(h.fileSystem)
	destination := getVersionsDir(h.fileSystem)

	return h.unzip.ExtractTarSource(destination, source)
}

func (h Helper) RenameGoDirectory(goVersionName string) error {
	versionDirName, err := h.GetLatestCreatedGoVersionDirectory()
	if err != nil {
		return err
	}

	target := fmt.Sprintf("%s/%s", getVersionsDir(h.fileSystem), versionDirName)

	if err := h.fileSystem.Rename(target, fmt.Sprintf("%s/%s", getVersionsDir(h.fileSystem), goVersionName)); err != nil {
		return err
	}

	return nil
}

func (h Helper) RemoveTarFile() error {
	if err := h.fileSystem.Remove(getTarFile(h.fileSystem)); err != nil {
		return err
	}

	return nil
}

func (h Helper) CreateExecutableSymlink(goVersionName string) error {
	destination := fmt.Sprintf("%s/%s/bin", getVersionsDir(h.fileSystem), goVersionName)

	files, err := h.fileSystem.ReadDir(destination)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		newFile := fmt.Sprintf("%s/%s", destination, f.Name())
		link := fmt.Sprintf("%s/%s", getBinDir(h.fileSystem), f.Name())

		// remove the symlink if exists already
		if _, err := h.fileSystem.Lstat(link); err == nil {
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
	path := getCurrentVersionFile(h.fileSystem)

	file, err := h.fileSystem.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = h.fileSystem.WriteString(file, goVersionName)
	return err
}

func (h Helper) StoreVersionsResponse(body []byte) error {
	return h.fileSystem.WriteFile(fmt.Sprintf("%s/%s", getAppDir(h.fileSystem), versionResponseFile), body, 0644)
}

func (h Helper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	byte_versions, err := h.fileSystem.ReadFile(fmt.Sprintf("%s/%s", getAppDir(h.fileSystem), versionResponseFile))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(byte_versions, &v); err != nil {
		return err
	}

	return nil
}

func (h Helper) AreVersionsCached() bool {
	file := fmt.Sprintf("%s/%s", getAppDir(h.fileSystem), versionResponseFile)

	if info, err := h.fileSystem.Stat(file); err == nil {
		currentTime := time.Now()
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week
		return currentTime.Sub(info.ModTime()).Hours() < 24*7
	}

	return false
}

func (h Helper) GetRecentVersion() string {
	path := getCurrentVersionFile(h.fileSystem)

	content, err := h.fileSystem.ReadFile(path)
	if err != nil {
		h.log.Error(err.Error())
		return ""
	}
	return string(content)
}

func (h Helper) DirectoryExists(goVersion string) bool {
	target := getVersionsDir(h.fileSystem)

	if _, err := h.fileSystem.Stat(fmt.Sprintf("%s/%s", target, goVersion)); !os.IsNotExist(err) {
		return true
	}

	return false
}

func (h Helper) DeleteDirectory(dirName string) error {
	target := getVersionsDir(h.fileSystem)
	err := h.fileSystem.RemoveAll(fmt.Sprintf("%s/%s", target, dirName))
	if err != nil {
		return err
	}

	return nil
}

func (h Helper) CreateInitFiles() (*os.File, error) {
	if err := h.fileSystem.MkdirIfNotExist(getAppDir(h.fileSystem), 0755); err != nil {
		return nil, err
	}
	if err := h.fileSystem.MkdirIfNotExist(getVersionsDir(h.fileSystem), 0755); err != nil {
		return nil, err
	}
	if err := h.fileSystem.MkdirIfNotExist(getBinDir(h.fileSystem), 0755); err != nil {
		return nil, err
	}

	// create log file
	filename := fmt.Sprintf("%s/%s", getAppDir(h.fileSystem), logFile)
	return h.fileSystem.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}

func (h Helper) GetLatestCreatedGoVersionDirectory() (string, error) {
	dir := getVersionsDir(h.fileSystem)
	files, err := h.fileSystem.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var modTime time.Time

	var dirName string
	for _, file := range files {
		if file.Type().IsDir() {
			info, err := file.Info()
			if err != nil {
				return "", err
			}
			info.ModTime()

			if !info.ModTime().Before(modTime) {
				if info.ModTime().After(modTime) {
					modTime = info.ModTime()
				}
				dirName = file.Name()
			}
		}
	}

	return dirName, nil
}

func New(fs FS, unzipper unzip.Unzipper, log *logger.Log) *Helper {
	return &Helper{
		fileSystem: fs,
		unzip:      unzipper,
		log:        log,
	}
}
