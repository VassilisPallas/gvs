// Package files provides interfaces for reading
// and writing files.
package files

import (
	"fmt"
)

var (
	// appDir contains the main directory for the CLI to store any data the CLI needs.
	appDir = ".gvs"
	// versionResponseFile contains the file name where the the response from the fetch request is stored for later use.
	versionResponseFile = "goVersions.json"
	// tarFileName contains the file name for the downloaded tar file.
	tarFileName = "downloaded.tar.gz"
	// goVersionsDir contains the directory name where the downloaded versions are stored (as well as the CURRENT file).
	goVersionsDir = ".go.versions"
	// binDir contains the directory when the symlinks are created for the selected version.
	binDir = "bin"
	// currentVersionFileName contains the file name where the currect (used) Go version is stored
	currentVersionFileName = "CURRENT"
	// logFile contains the file name where the logs are stored for debugging.
	logFile = "gvs.log"
)

// getAppDir returns the path for the `.gvs/` directory.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getAppDir(fs FS) string {
	return fmt.Sprintf("%s/%s", fs.GetHomeDirectory(), appDir)
}

// getVersionsDir returns the path for the `.go.versions/` directory.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getVersionsDir(fs FS) string {
	return fmt.Sprintf("%s/%s", getAppDir(fs), goVersionsDir)
}

// getTarFile returns the path for the `downloaded.tar.gz` file.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getTarFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getVersionsDir(fs), tarFileName)
}

// getCurrentVersionFile returns the path for the `CURRENT` file.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getCurrentVersionFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getVersionsDir(fs), currentVersionFileName)
}

// getBinDir returns the path for the `bin/` directory.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getBinDir(fs FS) string {
	return fmt.Sprintf("%s/%s", fs.GetHomeDirectory(), binDir)
}

// getVersionsResponseFile returns the path for the `goVersions.json` file.
//
// It is using the FS interface to get the $HOME directory, which is used as the starting point.
func getVersionsResponseFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getAppDir(fs), versionResponseFile)
}
