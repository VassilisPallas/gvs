// Package files provides interfaces for reading
// and writing files.
package files

import (
	"fmt"
)

var (
	appDir                 = ".gvs"
	versionResponseFile    = "goVersions.json"
	tarFileName            = "downloaded.tar.gz"
	goVersionsDir          = ".go.versions"
	binDir                 = "bin"
	currentVersionFileName = "CURRENT"
	logFile                = "gvs.log"
)

func getAppDir(fs FS) string {
	return fmt.Sprintf("%s/%s", fs.GetHomeDirectory(), appDir)
}

func getVersionsDir(fs FS) string {
	return fmt.Sprintf("%s/%s", getAppDir(fs), goVersionsDir)
}

func getTarFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getVersionsDir(fs), tarFileName)
}

func getCurrentVersionFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getVersionsDir(fs), currentVersionFileName)
}

func getBinDir(fs FS) string {
	return fmt.Sprintf("%s/%s", fs.GetHomeDirectory(), binDir)
}

func getVersionsResponseFile(fs FS) string {
	return fmt.Sprintf("%s/%s", getAppDir(fs), versionResponseFile)
}
