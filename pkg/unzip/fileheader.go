// Package unzip provides an interface for extracting
// tar and zip files.
package unzip

import "io/fs"

type FileHeader struct {
	// Name is the name of the file.
	//
	// It must be a relative path, not start with a drive letter (such as "C:"),
	// and must use forward slashes instead of back slashes. A trailing slash
	// indicates that this file is a directory and should have no data.
	Name string

	// A FileInfo describes a file and is returned by Stat.
	Info fs.FileInfo
}

// FileInfo returns an fs.FileInfo for the FileHeader.
func (h *FileHeader) FileInfo() fs.FileInfo {
	return h.Info
}
