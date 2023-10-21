package testutils

import (
	ioFS "io/fs"
)

type FakeDirEntry struct {
	DirEntryName  string
	DirEntryIsDir bool
	DireEntryType ioFS.FileMode

	DirEntryInfo      ioFS.FileInfo
	DirEntryInfoError error
}

func (de FakeDirEntry) Name() string {
	return de.DirEntryName
}

func (de FakeDirEntry) IsDir() bool {
	return de.DirEntryIsDir
}

func (de FakeDirEntry) Type() ioFS.FileMode {
	return de.DireEntryType
}

func (de FakeDirEntry) Info() (ioFS.FileInfo, error) {
	return de.DirEntryInfo, de.DirEntryInfoError
}
