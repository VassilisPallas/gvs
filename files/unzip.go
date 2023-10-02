package files

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Unziper interface {
	UnzipSource(dst string, src string) error
	UnzipFile(tarReader *tar.Reader, header *tar.Header, dst string) error
}

type Unzip struct {
	fs FS
}

func (u Unzip) UnzipSource(dst string, src string) error {
	reader, err := u.fs.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

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

		err = u.UnzipFile(tarReader, header, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u Unzip) UnzipFile(tarReader *tar.Reader, header *tar.Header, dst string) error {
	path := filepath.Join(dst, header.Name)
	info := header.FileInfo()

	// Check if file paths are not vulnerable to Zip Slip
	if !strings.HasPrefix(path, filepath.Clean(dst)+string(os.PathSeparator)) {
		// TODO: make custom error
		return fmt.Errorf("invalid file path: %s", path)
	}

	if info.IsDir() {
		if err := u.fs.MkdirAll(path, info.Mode()); err != nil {
			return err
		}

		return nil
	}

	if err := u.fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := u.fs.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := u.fs.Copy(file, tarReader); err != nil {
		return err
	}

	return nil
}
