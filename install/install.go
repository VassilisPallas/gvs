// Package install provides an interface for installing
// Go versions
package install

import (
	"context"
	"io"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/errors"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/logger"
)

// Installer is the interface that wraps the basic methods for installing new or existing versions.
// And existing version is a version that has been downloaded in the past and therefore the contents
// still exist on the local sysem.
//
// UnzipSource unzips the file path that defined as the source to the destination path.
// UnzipSource must return a non-null error if the unzip fails.
type Installer interface {
	NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error
	ExistingVersion(goVersionName string) error
}

// Install is the struct that implements the Installer interface
//
// Install structs accepts three fields, the fileHelpers that will be used to access and write files,
// ,the clientAPI that will be used to download the selected version (if it's a new version), and the
// log which is the logger.
type Install struct {
	fileHelpers files.FileHelpers
	clientAPI   api_client.GoClientAPI
	log         *logger.Log
}

// compareChecksums compares the SHA256 Checksum from the downloaded file with the checksum that was recieved from the API call
// that fetches all the versions and the infromation for each one of them.
//
// If there is a mismatch, compareChecksums will return an error.
func (i Install) compareChecksums(checksum string) error {
	hash, err := i.fileHelpers.GetTarChecksum()
	if err != nil {
		removeTarErr := i.fileHelpers.RemoveTarFile()
		// TODO: send to logger instead
		if removeTarErr != nil {
			return removeTarErr
		}

		return err
	}

	if hash != checksum {
		return &errors.ChecksumMisMatchError{Checksum: checksum, Hash: hash}
	}

	return nil
}

// createSymlink creates the symbolik links and updates the file that holds the currently installed version.
//
// If any of the above operations fail, createSymlink will return an error.
func (i Install) createSymlink(goVersionName string) error {
	if err := i.fileHelpers.CreateExecutableSymlink(goVersionName); err != nil {
		return err
	}

	if err := i.fileHelpers.UpdateRecentVersion(goVersionName); err != nil {
		return err
	}

	return nil
}

// newVersionHandler is the request callback that handles all the logic to install the new version.
//
// newVersionHandler creates the tar file from the response body and then after validating the checksum,
// it is unzipping the file and creates the symbolink liks. The unzipped directory is renamed to the selected
// Go version. For example if the version name is `1.20.7`, the directory that contains the unzipped files will also
// be named `1.20.7`.
//
// If any of the above operations fail, newVersionHandler will return an error.
func (i Install) newVersionHandler(checksum string, goVersionName string) func(content io.ReadCloser) error {
	return func(content io.ReadCloser) error {
		if err := i.fileHelpers.CreateTarFile(content); err != nil {
			return err
		}

		i.log.PrintMessage("Compare Checksums...\n")
		if err := i.compareChecksums(checksum); err != nil {
			return err
		}

		i.log.PrintMessage("Unzipping...\n")
		if err := i.fileHelpers.UnzipTarFile(); err != nil {
			return err
		}

		if err := i.fileHelpers.RenameGoDirectory(goVersionName); err != nil {
			return err
		}

		if err := i.fileHelpers.RemoveTarFile(); err != nil {
			return err
		}

		i.log.PrintMessage("Installing version...\n")
		return i.createSymlink(goVersionName)
	}
}

// NewVersion installs downloads and installs the selected version.
//
// NewVersion is first making a request to download the tar file (using the clientAPI interface),
// where it also passes the expected callback to handle the new version install logic.
//
// If the request or the version install fails, NewVersion will return an error.
func (i Install) NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error {
	i.log.PrintMessage("Downloading...\n")
	return i.clientAPI.DownloadVersion(ctx, fileName, i.newVersionHandler(checksum, goVersionName))
}

// ExistingVersion installs again an already existing version as the current go version.
// And existing version is a version that has been downloaded in the past and therefore the contents
// still exist on the local sysem.
//
// Since the version is already downloaded before, ExistingVersion only creates the new symblink link
// to override the existing go version.
//
// If the symbolink link creating fails, ExistingVersion will return an error.
func (i Install) ExistingVersion(goVersionName string) error {
	i.log.PrintMessage("Installing version...\n")
	return i.createSymlink(goVersionName)
}

// New returns a Install instance that implements the Installer interface.
// Each call to New returns a distinct Install instance even if the parameters are identical.
func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI, logger *logger.Log) Install {
	return Install{
		fileHelpers: fileHelpers,
		clientAPI:   clientAPI,
		log:         logger,
	}
}
