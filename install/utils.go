package install

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"

	"log"
	"os"
	"path/filepath"

	"github.com/VassilisPallas/gvs/files"
)

type InstallHelper interface {
	CreateTarFile(content io.ReadCloser) error
	GetTarChecksum() (string, error)
	UnzipTarFile() error
	RenameGoDirectory(goVersionName string) error
	RemoveTarFile() error
	CreateExecutableSymlink(goVersionName string) error
	UpdateRecentVersion(goVersionName string) error
}

type Helper struct {
	FileUtils files.FileUtils

	InstallHelper
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

	// TODO: add test
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

	// TODO: add test
	_, err = io.WriteString(file, goVersionName)
	return err
}
