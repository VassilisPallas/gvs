package files

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	FileUtils FileUtils

	FileHelpers
}

func (h Helper) CreateTarFile(content io.ReadCloser) error {
	file, err := os.Create(h.FileUtils.GetTarFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func (h Helper) GetTarChecksum() (string, error) {
	hasher := sha256.New()
	path := h.FileUtils.GetTarFile()

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// TODO: add tests
func (h Helper) UnzipTarFile() error {
	target := h.FileUtils.GetVersionsDir()

	reader, err := os.Open(h.FileUtils.GetTarFile())
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}

			continue
		} else {
			// This is happening only on 1.21.0. The directories are not
			// able to be found and instead the files contain the whole path
			// instead of just their name. This creates any missing directories
			// that exist on the file path. The permissions also are updated to
			// match the directory permissions from the previous versions
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	var deferError error
	defer func() {
		if err := reader.Close(); err != nil {
			deferError = err
		}

		if err := uncompressedStream.Close(); err != nil {
			deferError = err
		}
	}()

	return deferError
}

func (h Helper) RenameGoDirectory(goVersionName string) error {
	target := fmt.Sprintf("%s/%s", h.FileUtils.GetVersionsDir(), "go")

	if err := os.Rename(target, fmt.Sprintf("%s/%s", h.FileUtils.GetVersionsDir(), goVersionName)); err != nil {
		return err
	}

	return nil
}

func (h Helper) RemoveTarFile() error {
	err := os.Remove(h.FileUtils.GetTarFile())
	if err != nil {
		return err
	}
	return nil
}

// TODO: add tests
func (h Helper) CreateExecutableSymlink(goVersionName string) error {
	target := fmt.Sprintf("%s/%s/bin", h.FileUtils.GetVersionsDir(), goVersionName)

	files, err := os.ReadDir(target)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		newFile := fmt.Sprintf("%s/%s", target, f.Name())
		link := fmt.Sprintf("%s/%s", h.FileUtils.GetBinDir(), f.Name())

		if _, err := os.Lstat(link); err == nil {
			os.Remove(link)
		}

		if err := os.Symlink(newFile, link); err != nil {
			return err
		}
		os.Chmod(link, 0700)
	}

	return nil
}

func (h Helper) UpdateRecentVersion(goVersionName string) error {
	path := h.FileUtils.GetCurrentVersionFile()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, goVersionName)
	return err
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

func New(fileUtils FileUtils) *Helper {
	return &Helper{
		FileUtils: fileUtils,
	}
}
