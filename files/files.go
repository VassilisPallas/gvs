// Package files provides interfaces for reading
// and writing files.
package files

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/clock"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/pkg/unzip"
)

// FileHelpers is the interface that wraps the basic methods for working with files.
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

// Helper is the struct that implements the FileHelpers interface
type Helper struct {
	// fileSystem is the interface that is used as a wrapper for reading and writing files to the system.
	fileSystem FS
	// unzip is the interface that is used to unzip the downloaded version files
	unzip unzip.Unzipper
	// clock is the interface for time and duration.
	clock clock.Clock
	// log is the custom Logger
	log *logger.Log
}

// CreateTarFile creates the archive file based on the response from the API call
// that returns the file binary.
//
// If the creation of the file fails, CreateTarFile will return an error.
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

// GetTarChecksum returns the checksum for the downloaded the archive file.
//
// If for any reason if fails, GetTarChecksum returns back an empty checksum and the error.
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

// UnzipTarFile extracts the downloaded archive file.
//
// UnzipTarFile is using the unzip interface.
//
// If for any reason if fails, UnzipTarFile returns back an error.
func (h Helper) UnzipTarFile() error {
	source := getTarFile(h.fileSystem)
	destination := getVersionsDir(h.fileSystem)

	return h.unzip.ExtractTarSource(destination, source)
}

// RenameGoDirectory renames the extracted directory to the version name.
//
// After extracting the archive file, we need to rename it to the version name (e.g. go1.21.3).
// Since we can have multiple versions downloaded, this helps to store the versions into the
// correct directories.
//
// If for any reason if fails, RenameGoDirectory returns back an error.
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

// RemoveTarFile removes the archive file.
//
// If for any reason if fails, RemoveTarFile returns back an error.
func (h Helper) RemoveTarFile() error {
	if err := h.fileSystem.Remove(getTarFile(h.fileSystem)); err != nil {
		return err
	}

	return nil
}

// CreateExecutableSymlink creates the symlinks in $HOME/bin directory.
//
// The symlinks that will be created are the ones that can be found in
// the bin directory inside the version directory (e.g. ~/.gvs/.go.versions/go1.21.3/bin).
// If the symlinks exist already, we first remove the existing ones, and then creare the new ones.
//
// If for any reason if fails, CreateExecutableSymlink returns back an error.
func (h Helper) CreateExecutableSymlink(goVersionName string) error {
	versionBinDirectory := fmt.Sprintf("%s/%s/bin", getVersionsDir(h.fileSystem), goVersionName)

	files, err := h.fileSystem.ReadDir(versionBinDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		newFile := fmt.Sprintf("%s/%s", versionBinDirectory, file.Name())
		link := fmt.Sprintf("%s/%s", getBinDir(h.fileSystem), file.Name())

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

// UpdateRecentVersion updated the ~/.gvs/.go.versions/CURRENT file with the new installed version.
//
// We store the new installed version in this file, so we know which is the current used version.
//
// If for any reason if fails, UpdateRecentVersion returns back an error.
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

// StoreVersionsResponse stores the response from the fetch request so it can be used by avoid
// multiple requests to the API.
//
// If for any reason if fails, StoreVersionsResponse returns back an error.
func (h Helper) StoreVersionsResponse(body []byte) error {
	return h.fileSystem.WriteFile(getVersionsResponseFile(h.fileSystem), body, 0644)
}

// GetCachedResponse returns the cached response from fetch request.
//
// It will parse the JSON-encoded data and store it in the value pointed to by v.
//
// If for any reason if fails, GetCachedResponse returns back an error.
func (h Helper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	byte_versions, err := h.fileSystem.ReadFile(getVersionsResponseFile(h.fileSystem))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(byte_versions, &v); err != nil {
		return err
	}

	return nil
}

// AreVersionsCached returns if either the response from the fetch request is already
// cached or not.
//
// In case the cached file is older than a week, it will return false to force
// the caching again.
func (h Helper) AreVersionsCached() bool {
	file := getVersionsResponseFile(h.fileSystem)

	if info, err := h.fileSystem.Stat(file); err == nil {
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week.
		// This helps to purge the cached file and return fresh data.
		return h.clock.GetDiffInHoursFromNow(info.ModTime()) < 24*7
	}

	return false
}

// GetRecentVersion returns the current (installed) Go version from ~/.gvs/.go.versions/CURRENT
//
// If for any reason if fails, GetRecentVersion logs the error message and returns back an empty string.
func (h Helper) GetRecentVersion() string {
	path := getCurrentVersionFile(h.fileSystem)

	content, err := h.fileSystem.ReadFile(path)
	if err != nil {
		h.log.Error(err.Error())
		return ""
	}
	return string(content)
}

// DirectoryExists checks if the given Go version directory exists or not.
func (h Helper) DirectoryExists(goVersion string) bool {
	target := getVersionsDir(h.fileSystem)

	_, err := h.fileSystem.Stat(fmt.Sprintf("%s/%s", target, goVersion))
	return err == nil
}

// DeleteDirectory deletes the given Go version directory.
//
// If for any reason if fails, DeleteDirectory returns back an error.
func (h Helper) DeleteDirectory(goVersion string) error {
	target := getVersionsDir(h.fileSystem)
	return h.fileSystem.RemoveAll(fmt.Sprintf("%s/%s", target, goVersion))
}

// CreateInitFiles creates the files that are required for the CLI.
//
// It is creating the below:
// `$HOME/.gvs/` - The main directory for the CLI to store any data the CLI needs.
//
// `$HOME/.gvs/.go.versions/` - The directory where the downloaded versions are stored (as well as the CURRENT file).
//
// `$HOME/bin` - This is the directory when the symlinks are created for the selected version.
//
// `$HOME/.gvs/gvs.log` - This is where the logs are stored for debugging.
//
// Once the gvs.log is created, the *os.File is returned back so it can be used from the logger.
//
// If for any reason if fails, CreateInitFiles returns back nil for the file and the error.
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

// GetLatestCreatedGoVersionDirectory returns the name of the latest modified directory in .gvs/.go.versions/.
//
// If for any reason if fails, GetCachedResponse returns back an empty string for the name and the error.
func (h Helper) GetLatestCreatedGoVersionDirectory() (string, error) {
	dir := getVersionsDir(h.fileSystem)
	files, err := h.fileSystem.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var modTime time.Time

	var dirName string
	for _, file := range files {
		if file.IsDir() {
			info, err := file.Info()
			if err != nil {
				return "", err
			}

			info.ModTime()
			if !h.clock.IsBefore(info.ModTime(), modTime) {
				if h.clock.IsAfter(info.ModTime(), modTime) {
					modTime = info.ModTime()
				}
				dirName = file.Name()
			}
		}
	}

	return dirName, nil
}

// New returns a *Helper instance that implements the FileHelpers interface.
// Each call to New returns a distinct *Helper instance even if the parameters are identical.
func New(fs FS, clock clock.Clock, unzipper unzip.Unzipper, log *logger.Log) *Helper {
	return &Helper{
		fileSystem: fs,
		clock:      clock,
		unzip:      unzipper,
		log:        log,
	}
}
