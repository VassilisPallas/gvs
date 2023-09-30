package install

import (
	"context"
	"io"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/errors"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/logger"
)

type Installer interface {
	NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error
	ExistingVersion(goVersionName string) error
}

type Install struct {
	FileHelpers files.FileHelpers
	ClientAPI   api_client.GoClientAPI
	Log         *logger.Log

	Installer
}

func (i Install) compareChecksums(checksum string) error {
	hash, err := i.FileHelpers.GetTarChecksum()
	if err != nil {
		// TODO: test this
		i.FileHelpers.RemoveTarFile()
		return err
	}

	if hash != checksum {
		return &errors.ChecksumMisMatchError{Checksum: checksum, Hash: hash}
	}

	return nil
}

func (i Install) createSymlink(goVersionName string) error {
	if err := i.FileHelpers.CreateExecutableSymlink(goVersionName); err != nil {
		return err
	}

	if err := i.FileHelpers.UpdateRecentVersion(goVersionName); err != nil {
		return err
	}

	return nil
}

func (i Install) newVersionHandler(checksum string, goVersionName string) func(content io.ReadCloser) error {
	return func(content io.ReadCloser) error {
		if err := i.FileHelpers.CreateTarFile(content); err != nil {
			return err
		}

		i.Log.PrintMessage("Compare Checksums...\n")
		if err := i.compareChecksums(checksum); err != nil {
			return err
		}

		i.Log.PrintMessage("Unzipping...\n")
		if err := i.FileHelpers.UnzipTarFile(); err != nil {
			return err
		}

		if err := i.FileHelpers.RenameGoDirectory(goVersionName); err != nil {
			return err
		}

		if err := i.FileHelpers.RemoveTarFile(); err != nil {
			return err
		}

		i.Log.PrintMessage("Installing version...\n")
		return i.createSymlink(goVersionName)
	}
}

func (i Install) NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error {
	i.Log.PrintMessage("Downloading...\n")
	return i.ClientAPI.DownloadVersion(ctx, fileName, i.newVersionHandler(checksum, goVersionName))
}

func (i Install) ExistingVersion(goVersionName string) error {
	i.Log.PrintMessage("Installing version...\n")
	return i.createSymlink(goVersionName)
}

// TODO: check if should return *Install
func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI, logger *logger.Log) Install {
	return Install{FileHelpers: fileHelpers, ClientAPI: clientAPI, Log: logger}
}
