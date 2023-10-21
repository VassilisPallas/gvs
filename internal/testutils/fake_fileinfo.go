package testutils

import (
	ioFS "io/fs"
	"time"
)

type FakeFileInfo struct {
	FileName    string
	FileSize    int64
	FileMode    ioFS.FileMode
	FileModTime time.Time
	FileIsDir   bool
}

func (fi FakeFileInfo) Name() string {
	return fi.FileName
}

func (fi FakeFileInfo) Size() int64 {
	return fi.FileSize
}

func (fi FakeFileInfo) Mode() ioFS.FileMode {
	return fi.FileMode
}

func (fi FakeFileInfo) ModTime() time.Time {
	return fi.FileModTime
}

func (fi FakeFileInfo) IsDir() bool {
	return fi.FileIsDir
}

func (FakeFileInfo) Sys() any {
	return nil
}
