package testutils

import (
	"io"
	ioFS "io/fs"
	"os"
)

type FakeFileSystem struct {
	HomeDir string

	CreateMockFile *os.File
	CreateError    error

	CopyError error

	OpenMockFile *os.File
	OpenError    error
}

func (FakeFileSystem) Chmod(name string, mode ioFS.FileMode) error {
	return nil
}

func (fs FakeFileSystem) Create(name string) (*os.File, error) {
	return fs.CreateMockFile, fs.CreateError
}

func (fs FakeFileSystem) Open(name string) (*os.File, error) {
	return fs.OpenMockFile, fs.OpenError
}

func (FakeFileSystem) OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error) {
	return nil, nil
}

func (FakeFileSystem) ReadDir(name string) ([]ioFS.DirEntry, error) {
	return nil, nil
}

func (FakeFileSystem) ReadFile(name string) ([]byte, error) {
	return nil, nil
}

func (FakeFileSystem) Stat(name string) (ioFS.FileInfo, error) {
	return nil, nil
}

func (FakeFileSystem) Lstat(name string) (ioFS.FileInfo, error) {
	return nil, nil
}

func (fs FakeFileSystem) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return 0, fs.CopyError
}

func (FakeFileSystem) Symlink(oldname string, newname string) error {
	return nil
}

func (FakeFileSystem) WriteFile(name string, data []byte, perm ioFS.FileMode) error {
	return nil
}

func (FakeFileSystem) WriteString(w io.Writer, s string) (n int, err error) {
	return 0, nil
}

func (FakeFileSystem) MkdirAll(path string, perm ioFS.FileMode) error {
	return nil
}

func (FakeFileSystem) MkdirIfNotExist(path string, perm ioFS.FileMode) error {
	return nil
}

func (FakeFileSystem) Rename(oldpath string, newpath string) error {
	return nil
}

func (FakeFileSystem) Remove(name string) error {
	return nil
}

func (FakeFileSystem) RemoveAll(path string) error {
	return nil
}

func (fs FakeFileSystem) GetHomeDirectory() string {
	return fs.HomeDir
}
