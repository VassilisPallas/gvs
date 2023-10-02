package files

import (
	"fmt"
	"os"
	"os/user"
)

type FileUtils interface {
	GetHomeDirectory() string
	GetAppDir() string
	GetVersionsDir() string
	GetTarFile() string
	GetBinDir() string
	GetCurrentVersionFile() string
	GetVersionResponseFile() string
	CreateInitFiles() error
	CreateLogFile() (*os.File, error)
}

type Files struct {
	appDir                 string
	versionResponseFile    string
	tarFileName            string
	goVersionsDir          string
	binDir                 string
	currentVersionFileName string
	logFile                string

	FileSystem FS

	FileUtils
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

// TODO: add tests
func (f Files) CreateInitFiles() error {
	if err := f.FileSystem.MkdirIfNotExist(f.GetAppDir(), 0755); err != nil {
		return err
	}
	if err := f.FileSystem.MkdirIfNotExist(f.GetVersionsDir(), 0755); err != nil {
		return err
	}
	if err := f.FileSystem.MkdirIfNotExist(f.GetBinDir(), 0755); err != nil {
		return err
	}

	return nil
}

// TODO: add tests
func (f Files) CreateLogFile() (*os.File, error) {
	filename := fmt.Sprintf("%s/%s", f.GetAppDir(), f.logFile)
	return f.FileSystem.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
		FileSystem:             FileSystem{},
	}
}
