// Package files provides interfaces for reading
// and writing files.
package files

import (
	"io"
	ioFS "io/fs"
	"os"
)

// FS is the interface that wraps the basic methods for reading and writing files to the system.
type FS interface {
	// Chmod changes the mode of the named file to mode.
	Chmod(name string, mode ioFS.FileMode) error
	// Create creates or truncates the named file.
	Create(name string) (*os.File, error)

	// Open opens the named file for reading.
	Open(name string) (*os.File, error)
	// OpenFile opens the named file with specified flag.
	OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error)
	// ReadDir reads the named directory, returning all its directory entries sorted by filename.
	ReadDir(name string) ([]ioFS.DirEntry, error)
	// ReadFile reads the named file and returns the contents.
	ReadFile(name string) ([]byte, error)
	// Stat returns a FileInfo describing the named file.
	Stat(name string) (ioFS.FileInfo, error)
	// Lstat returns a FileInfo describing the named file. If the file is a symbolic link, the returned FileInfo describes the symbolic link.
	Lstat(name string) (ioFS.FileInfo, error)

	// Copy copies from src to dst until either EOF is reached on src or an error occurs. It returns the number of bytes copied and the first error encountered while copying, if any.
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	// Symlink creates newname as a symbolic link to oldname.
	Symlink(oldname string, newname string) error

	// WriteFile writes data to the named file, creating it if necessary.
	WriteFile(name string, data []byte, perm ioFS.FileMode) error
	// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
	WriteString(w io.Writer, s string) (n int, err error)

	// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
	MkdirAll(path string, perm ioFS.FileMode) error
	// MkdirIfNotExist creates a directory named path if it does not exist, along with any necessary parents, and returns nil, or else returns an error.
	MkdirIfNotExist(path string, perm ioFS.FileMode) error

	// Rename renames (moves) oldpath to newpath.
	Rename(oldpath string, newpath string) error

	// Remove removes the named file or (empty) directory.
	Remove(name string) error
	// RemoveAll removes path and any children it contains.
	RemoveAll(path string) error

	// GetHomeDirectory return back the home directory for the current user.
	GetHomeDirectory() string
}

// FileSystem is the struct that implements the FS interface
type FileSystem struct{}

// Chmod changes the mode of the named file to mode.
//
// It is a wrapper for the os.Chmod function.
func (FileSystem) Chmod(name string, mode ioFS.FileMode) error {
	return os.Chmod(name, mode)
}

// Create creates or truncates the named file.
//
// It is a wrapper for the os.Create function.
func (FileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Open opens the named file for reading.
//
// It is a wrapper for the os.Open function.
func (FileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

// OpenFile opens the named file with specified flag.
//
// It is a wrapper for the os.OpenFile function.
func (FileSystem) OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// ReadDir reads the named directory, returning all its directory entries sorted by filename.
//
// It is a wrapper for the os.ReadDir function.
func (FileSystem) ReadDir(name string) ([]ioFS.DirEntry, error) {
	return os.ReadDir(name)
}

// ReadFile reads the named file and returns the contents.
//
// It is a wrapper for the os.ReadFile function.
func (FileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Stat returns a FileInfo describing the named file.
//
// It is a wrapper for the os.Stat function.
func (FileSystem) Stat(name string) (ioFS.FileInfo, error) {
	return os.Stat(name)
}

// Lstat returns a FileInfo describing the named file. If the file is a symbolic link, the returned FileInfo describes the symbolic link.
//
// It is a wrapper for the os.Lstat function.
func (FileSystem) Lstat(name string) (ioFS.FileInfo, error) {
	return os.Lstat(name)
}

// Copy copies from src to dst until either EOF is reached on src or an error occurs. It returns the number of bytes copied and the first error encountered while copying, if any.
//
// It is a wrapper for the io.Copy function.
func (FileSystem) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

// Symlink creates newname as a symbolic link to oldname.
//
// It is a wrapper for the os.Symlink function.
func (FileSystem) Symlink(oldname string, newname string) error {
	return os.Symlink(oldname, newname)
}

// WriteFile writes data to the named file, creating it if necessary.
//
// It is a wrapper for the os.WriteFile function.
func (FileSystem) WriteFile(name string, data []byte, perm ioFS.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
//
// It is a wrapper for the io.WriteString function.
func (FileSystem) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
//
// It is a wrapper for the os.MkdirAll function.
func (FileSystem) MkdirAll(path string, perm ioFS.FileMode) error {
	return os.MkdirAll(path, perm)
}

// MkdirIfNotExist creates a directory named path if it does not exist, along with any necessary parents, and returns nil, or else returns an error.
//
// It is using the Stat method from the FS interface to check if the directory exists already.
// Then is using the MkdirAll method from the FS interface to create the directory, along with any necessary parents.
//
// It returns an error either when Stat fails or when the directory creation fails.
func (fs FileSystem) MkdirIfNotExist(path string, perm ioFS.FileMode) error {
	_, err := fs.Stat(path)

	if os.IsNotExist(err) {
		err = fs.MkdirAll(path, perm)
		if err != nil {
			return err
		}
	}
	return err
}

// Rename renames (moves) oldpath to newpath.
//
// It is a wrapper for the os.Rename function.
func (FileSystem) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// Remove removes the named file or (empty) directory.
//
// It is a wrapper for the os.Remove function.
func (FileSystem) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll removes path and any children it contains.
//
// It is a wrapper for the os.RemoveAll function.
func (FileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// GetHomeDirectory return back the home directory for the current user.
//
// If an error occurs, GetHomeDirectory will panic.
func (FileSystem) GetHomeDirectory() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return dirname
}
