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
}

type Helper struct {
	FileUtils  FileUtils
	Unzip      Unziper
	FileSystem FS

	FileHelpers
}

func (h Helper) CreateTarFile(content io.ReadCloser) error {
	file, err := h.FileSystem.Create(h.FileUtils.GetTarFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = h.FileSystem.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func (h Helper) GetTarChecksum() (string, error) {
	hasher := sha256.New()
	path := h.FileUtils.GetTarFile()

	f, err := h.FileSystem.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := h.FileSystem.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// TODO: add tests
func (h Helper) UnzipTarFile() error {
	source := h.FileUtils.GetTarFile()
	destination := h.FileUtils.GetVersionsDir()

	return h.Unzip.UnzipSource(destination, source)
}

func (h Helper) RenameGoDirectory(goVersionName string) error {
	target := fmt.Sprintf("%s/%s", h.FileUtils.GetVersionsDir(), "go")

	if err := h.FileSystem.Rename(target, fmt.Sprintf("%s/%s", h.FileUtils.GetVersionsDir(), goVersionName)); err != nil {
		return err
	}

	return nil
}

func (h Helper) RemoveTarFile() error {
	err := h.FileSystem.Remove(h.FileUtils.GetTarFile())
	if err != nil {
		return err
	}
	return nil
}

// TODO: add tests
func (h Helper) CreateExecutableSymlink(goVersionName string) error {
	destination := fmt.Sprintf("%s/%s/bin", h.FileUtils.GetVersionsDir(), goVersionName)

	files, err := h.FileSystem.ReadDir(destination)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		newFile := fmt.Sprintf("%s/%s", destination, f.Name())
		link := fmt.Sprintf("%s/%s", h.FileUtils.GetBinDir(), f.Name())

		if _, err := h.FileSystem.Lstat(link); err == nil {
			h.FileSystem.Remove(link)
		}

		if err := h.FileSystem.Symlink(newFile, link); err != nil {
			return err
		}
		h.FileSystem.Chmod(link, 0700)
	}

	return nil
}

func (h Helper) UpdateRecentVersion(goVersionName string) error {
	path := h.FileUtils.GetCurrentVersionFile()

	file, err := h.FileSystem.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, goVersionName)
	return err
}

func (h Helper) StoreVersionsResponse(body []byte) error {
	return h.FileSystem.WriteFile(fmt.Sprintf("%s/%s", h.FileUtils.GetAppDir(), h.FileUtils.GetVersionResponseFile()), body, 0644)
}

func (h Helper) GetCachedResponse(v *[]api_client.VersionInfo) error {
	byte_versions, err := h.FileSystem.ReadFile(fmt.Sprintf("%s/%s", h.FileUtils.GetAppDir(), h.FileUtils.GetVersionResponseFile()))
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

	if info, err := h.FileSystem.Stat(file); err == nil {
		currentTime := time.Now()
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week
		return currentTime.Sub(info.ModTime()).Hours() < 24*7
	}

	return false
}

func (h Helper) GetRecentVersion() string {
	path := h.FileUtils.GetCurrentVersionFile()

	content, err := h.FileSystem.ReadFile(path)
	if err != nil {
		// TODO: pass error as log to file
		return ""
	}
	return string(content)
}

func (h Helper) DirectoryExists(goVersion string) bool {
	target := h.FileUtils.GetVersionsDir()

	if _, err := h.FileSystem.Stat(fmt.Sprintf("%s/%s", target, goVersion)); !os.IsNotExist(err) {
		return true
	}

	return false
}

func (h Helper) DeleteDirectory(dirName string) error {
	target := h.FileUtils.GetVersionsDir()
	err := h.FileSystem.RemoveAll(fmt.Sprintf("%s/%s", target, dirName))
	if err != nil {
		return err
	}

	return nil
}

func New(fileUtils FileUtils) *Helper {
	fs := FileSystem{}

	return &Helper{
		FileUtils:  fileUtils,
		Unzip:      Unzip{fs: fs},
		FileSystem: fs,
	}
}
