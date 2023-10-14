// Package files provides interfaces for reading
// and writing files.
package files

import (
	"fmt"
	"os/user"
)

// FileUtils is the interface that wraps the basic methods for making requests.
//
// FetchVersions fetches and returns the available Go versions.
// The versions should be parsed and stored in the value pointed to by v.
// FetchVersions must return a non-null error if the request, or the parsing of the response fails.
//
// DownloadVersion downloads the content (most likely a tar.gz file) and then is passing
// the response to the callack function.
// DownloadVersion must close the response body reader after passing it in the callback function.
// DownloadVersion must return an non-null error if the request failes or the callback returns
// an non-null error.
type FileUtils interface {
	GetHomeDirectory() string
	GetAppDir() string
	GetVersionsDir() string
	GetTarFile() string
	GetBinDir() string
	GetCurrentVersionFile() string
	GetVersionResponseFile() string
	GetLogFile() string
}

type Files struct {
	appDir                 string
	versionResponseFile    string
	tarFileName            string
	goVersionsDir          string
	binDir                 string
	currentVersionFileName string
	logFile                string
}

func (f Files) GetHomeDirectory() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.HomeDir
}

func (f Files) GetAppDir() string {
	return fmt.Sprintf("%s/%s", f.GetHomeDirectory(), f.appDir)
}

func (f Files) GetVersionsDir() string {
	return fmt.Sprintf("%s/%s", f.GetAppDir(), f.goVersionsDir)
}

func (f Files) GetTarFile() string {
	return fmt.Sprintf("%s/%s", f.GetVersionsDir(), f.tarFileName)
}

func (f Files) GetBinDir() string {
	return fmt.Sprintf("%s/%s", f.GetHomeDirectory(), f.binDir)
}

func (f Files) GetCurrentVersionFile() string {
	return fmt.Sprintf("%s/%s", f.GetVersionsDir(), f.currentVersionFileName)
}

func (f Files) GetVersionResponseFile() string {
	return f.versionResponseFile
}

func (f Files) GetLogFile() string {
	return f.logFile
}

func NewUtils() *Files {
	return &Files{
		appDir:                 ".gvs",
		versionResponseFile:    "goVersions.json",
		tarFileName:            "downloaded.tar.gz",
		goVersionsDir:          ".go.versions",
		binDir:                 "bin",
		currentVersionFileName: "CURRENT",
		logFile:                "gvs.log",
	}
}
