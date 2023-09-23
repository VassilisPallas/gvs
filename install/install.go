package install

import (
	"context"
	"fmt"
	"io"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/files"
)

type Installer interface {
	NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error
	ExistingVersion(goVersionName string) error
}

type Install struct {
	FileHelpers files.FileHelpers
	ClientAPI   api_client.GoClientAPI

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
		return fmt.Errorf("checksums do not match.\nExpected: %s\nGot: %s", checksum, hash)
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

		fmt.Println("Compare Checksums...")
		if err := i.compareChecksums(checksum); err != nil {
			return err
		}

		fmt.Println("Unzipping...")
		if err := i.FileHelpers.UnzipTarFile(); err != nil {
			return err
		}

		if err := i.FileHelpers.RenameGoDirectory(goVersionName); err != nil {
			return err
		}

		if err := i.FileHelpers.RemoveTarFile(); err != nil {
			return err
		}

		fmt.Println("Installing version...")
		return i.createSymlink(goVersionName)
	}
}

func (i Install) NewVersion(ctx context.Context, fileName string, checksum string, goVersionName string) error {
	fmt.Println("Downloading...")
	return i.ClientAPI.DownloadVersion(ctx, fileName, i.newVersionHandler(checksum, goVersionName))
}

func (i Install) ExistingVersion(goVersionName string) error {
	fmt.Println("Installing version...")
	return i.createSymlink(goVersionName)
}

func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI) Installer {
	return Install{FileHelpers: fileHelpers, ClientAPI: clientAPI}
}
