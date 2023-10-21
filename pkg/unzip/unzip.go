// Package unzip provides an interface for extracting
// tar and zip files.
package unzip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzipper is the interface that wraps the basic methods for extracting files.
type Unzipper interface {
	// ExtractTarSource extracts the file path that defined as the source (.tar.gz file) to the destination path.
	// ExtractTarSource must return a non-null *UnzipError if the extract fails.
	ExtractTarSource(dst string, src string) error

	// ExtractZipSource extracts the file path that defined as the source (.zip file) to the destination path.
	// ExtractZipSource must return a non-null *UnzipError if the extract fails.
	ExtractZipSource(dst string, src string) error
}

// Unzip is the struct that implements the Unzipper interface
type Unzip struct {
	// The FileSystem methods to be used for I/O and OS operations
	FileSystem FS
}

// unzipFile creates the given entry which can be either a directory or a file.
//
// For all the I/O and OS operation it is using the FS interface implementation.
//
// If the creation of the file or the directory fails, unzipFile returns back an not-null *UnzipError type of error.
func (u Unzip) unzipFile(src io.ReadCloser, header *FileHeader, dst string) error {
	path := filepath.Join(dst, header.Name)
	info := header.FileInfo()

	// Check if file paths are not vulnerable to Zip Slip
	if !strings.HasPrefix(path, filepath.Clean(dst)+string(os.PathSeparator)) {
		return &UnzipError{fmt.Errorf("invalid file path: %s", path)}
	}

	// create directory tree
	if info.IsDir() {
		if err := u.FileSystem.MkdirAll(path, info.Mode()); err != nil {
			return &UnzipError{err}
		}

		return nil
	}

	if err := u.FileSystem.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return &UnzipError{err}
	}

	file, err := u.FileSystem.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return &UnzipError{err}
	}
	defer file.Close()

	if _, err := u.FileSystem.Copy(file, src); err != nil {
		return &UnzipError{err}
	}
	defer src.Close()

	return nil
}

// ExtractTarSource extracts the file path that defined as the source (.tar.gz file) to the destination path.
//
// For all the I/O and OS operation it is using the FS interface implementation.
// For getting the stream reader, first it is using the NewReader from the gzip library
// and then the NewReader from the tar library, where the result of the first reader is passed
// as an input to the second reader.
//
// It then iterates on each entry in the tar archive and calls the unzipFile method to handle the file
// (or directory) creation.
//
// If any of the above operations fail, ExtractTarSource returns back an not-null *UnzipError type of error.
func (u Unzip) ExtractTarSource(dst string, src string) error {
	reader, err := u.FileSystem.Open(src)
	if err != nil {
		return &UnzipError{err}
	}
	defer reader.Close()

	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		return &UnzipError{err}
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if isEOF(err) {
			break
		}

		if err != nil {
			return &UnzipError{err}
		}

		fileHeader := &FileHeader{
			Name: header.Name,
			Info: header.FileInfo(),
		}

		unzipError := u.unzipFile(io.NopCloser(tarReader), fileHeader, dst)
		if unzipError != nil {
			return &UnzipError{err}
		}
	}

	return nil
}

// ExtractZipSource extracts the file path that defined as the source (.zip file) to the destination path.
//
// For all the I/O and OS operation it is using the FS interface implementation.
// For getting the stream reader, first it is using the OpenReader from the zip library.
//
// It then iterates on each file in the zip archive and calls the unzipFile method to handle the file
// (or directory) creation.
//
// If any of the above operations fail, ExtractZipSource returns back an not-null *UnzipError type of error.
func (u Unzip) ExtractZipSource(dst string, src string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return &UnzipError{err}
	}
	defer reader.Close()

	dst, err = filepath.Abs(dst)
	if err != nil {
		return &UnzipError{err}
	}

	for _, f := range reader.File {
		zippedFile, err := f.Open()
		if err != nil {
			return &UnzipError{err}
		}

		fileHeader := &FileHeader{
			Name: f.Name,
			Info: f.FileInfo(),
		}

		unzipError := u.unzipFile(zippedFile, fileHeader, dst)

		if unzipError != nil {
			return &UnzipError{err}
		}
	}

	return nil
}
