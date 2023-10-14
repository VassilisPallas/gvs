// Package files provides interfaces for reading
// and writing files.
package files

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unziper is the interface that wraps the basic methods for unzipping files.
//
// UnzipSource unzips the file path that defined as the source to the destination path.
// UnzipSource must return a non-null error if the unzip fails.
// TODO: add tests
type Unziper interface {
	UnzipSource(dst string, src string) error
}

// Unzip is the struct that implements the Unziper interface
//
// Unzip structs accepts the fs field, which is the wrapper for all the I/O and OS operations
// regarding reading and writing files.
type Unzip struct {
	fs FS
}

// unzipFile creates the given entry which can be either a directory or a file.
//
// For all the I/O and OS operation it is using the FS interface implementation.
//
// If the creation of the file or the directory fails, unzipFile returns back an not-null error.
func (u Unzip) unzipFile(tarReader *tar.Reader, header *tar.Header, dst string) error {
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

// UnzipSource unzips the file path that defined as the source to the destination path.
//
// For all the I/O and OS operation it is using the FS interface implementation.
// For getting the stream reader, first it is using the NewReader from the gzip library
// and then the NewReader from the tar library, where the result of the first reader is passed
// as an input to the second reader.
//
// It then iterates on each entry in the tar archive and calls the unzipFile from the same interface to
// handle the file (or directory) creation.
//
// If any of the above operations fail, UnzipSource returns back an not-null error.
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

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		err = u.unzipFile(tarReader, header, dst)
		if err != nil {
			return err
		}
	}

	return nil
}
