// Package unzip provides functions for extracting
// tar and zip files.
package unzip

import (
	"io"
	ioFS "io/fs"
	"os"
)

// FS is the interface that wraps the basic methods for reading and writing files to the system.
type FS interface {
	// Copy copies from src to dst until either EOF is reached on src or an error occurs. It returns the number of bytes copied and the first error encountered while copying, if any.
	Copy(dst io.Writer, src io.Reader) (written int64, err error)

	// Open opens the named file for reading.
	Open(name string) (*os.File, error)

	// OpenFile opens the named file with specified flag.
	OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error)

	// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
	MkdirAll(path string, perm ioFS.FileMode) error
}
