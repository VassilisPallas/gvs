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

	ReadDirMockResponse []ioFS.DirEntry
	ReadDirError        error

	RenameError error

	RemoveError error

	SymlinkError error

	ChmodError error

	WriteStringError error

	WriteFileError error

	ReadFileBytes []byte
	ReadFileError error

	StatMockResponse ioFS.FileInfo
	StatError        error

	RemoveAllError error

	MkdirIfNotExistPathToFail string
	MkdirIfNotExistError      error

	OpenFileMockFile *os.File
	OpenFileError    error
}

func (fs FakeFileSystem) Chmod(name string, mode ioFS.FileMode) error {
	return fs.ChmodError
}

func (fs FakeFileSystem) Create(name string) (*os.File, error) {
	return fs.CreateMockFile, fs.CreateError
}

func (fs FakeFileSystem) Open(name string) (*os.File, error) {
	return fs.OpenMockFile, fs.OpenError
}

func (fs FakeFileSystem) OpenFile(name string, flag int, perm ioFS.FileMode) (*os.File, error) {
	return fs.OpenFileMockFile, fs.OpenFileError
}

func (fs FakeFileSystem) ReadDir(name string) ([]ioFS.DirEntry, error) {
	return fs.ReadDirMockResponse, fs.ReadDirError
}

func (fs FakeFileSystem) ReadFile(name string) ([]byte, error) {
	return fs.ReadFileBytes, fs.ReadFileError
}

func (fs FakeFileSystem) Stat(name string) (ioFS.FileInfo, error) {
	return fs.StatMockResponse, fs.StatError
}

func (FakeFileSystem) Lstat(name string) (ioFS.FileInfo, error) {
	return nil, nil
}

func (fs FakeFileSystem) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return 0, fs.CopyError
}

func (fs FakeFileSystem) Symlink(oldname string, newname string) error {
	return fs.SymlinkError
}

func (fs FakeFileSystem) WriteFile(name string, data []byte, perm ioFS.FileMode) error {
	return fs.WriteFileError
}

func (fs FakeFileSystem) WriteString(w io.Writer, s string) (n int, err error) {
	return 0, fs.WriteStringError
}

func (FakeFileSystem) MkdirAll(path string, perm ioFS.FileMode) error {
	return nil
}

func (fs FakeFileSystem) MkdirIfNotExist(path string, perm ioFS.FileMode) error {
	if fs.MkdirIfNotExistPathToFail == path {
		return fs.MkdirIfNotExistError
	}

	return nil
}

func (fs FakeFileSystem) Rename(oldpath string, newpath string) error {
	return fs.RenameError
}

func (fs FakeFileSystem) Remove(name string) error {
	return fs.RemoveError
}

func (fs FakeFileSystem) RemoveAll(path string) error {
	return fs.RemoveAllError
}

func (fs FakeFileSystem) GetHomeDirectory() string {
	return fs.HomeDir
}
