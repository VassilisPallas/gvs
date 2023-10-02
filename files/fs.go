package files

import (
	"io"
	ioFS "io/fs"
	"os"
)

type FS interface {
	Chmod(name string, mode ioFS.FileMode) error
	Create(name string) (*os.File, error)

	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error)
	ReadDir(name string) ([]ioFS.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	Stat(name string) (ioFS.FileInfo, error)
	Lstat(name string) (ioFS.FileInfo, error)

	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	Symlink(oldname string, newname string) error
	WriteFile(name string, data []byte, perm ioFS.FileMode) error
	WriteString(w io.Writer, s string) (n int, err error)

	MkdirAll(path string, perm ioFS.FileMode) error
	MkdirIfNotExist(path string, perm ioFS.FileMode) error

	Rename(oldpath string, newpath string) error

	Remove(name string) error
	RemoveAll(path string) error
}

type FileSystem struct{}

func (FileSystem) Chmod(name string, mode ioFS.FileMode) error {
	return os.Chmod(name, mode)
}

func (FileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (FileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (FileSystem) OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (FileSystem) MkdirAll(path string, perm ioFS.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (FileSystem) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

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

func (FileSystem) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (FileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (FileSystem) ReadDir(name string) ([]ioFS.DirEntry, error) {
	return os.ReadDir(name)
}

func (FileSystem) Lstat(name string) (ioFS.FileInfo, error) {
	return os.Lstat(name)
}

func (FileSystem) Symlink(oldname string, newname string) error {
	return os.Symlink(oldname, newname)
}

func (FileSystem) WriteFile(name string, data []byte, perm ioFS.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (FileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (FileSystem) Stat(name string) (ioFS.FileInfo, error) {
	return os.Stat(name)
}

func (FileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (FileSystem) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}
