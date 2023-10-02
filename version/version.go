package version

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/errors"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/logger"
)

type Versioner interface {
	GetVersions(forceFetchVersions bool) ([]*ExtendedVersion, error)
	DeleteUnusedVersions(evs []*ExtendedVersion) (int, error)
	GetLatestVersion(evs []*ExtendedVersion) int
	Install(ev *ExtendedVersion, os string, arch string) error
	GetPromptVersions(evs []*ExtendedVersion, showAllVersions bool) []*ExtendedVersion

	FilterAlreadyDownloadedVersions(evs []*ExtendedVersion) []string
}

type Version struct {
	installer   install.Installer
	clientAPI   api_client.GoClientAPI
	fileHelpers files.FileHelpers
	log         *logger.Log
}

type ExtendedVersion struct {
	UsedVersion      bool
	AlreadyInstalled bool

	api_client.VersionInfo
}

func (ev *ExtendedVersion) addExtras(helper files.FileHelpers) {
	if helper.DirectoryExists(ev.Version) {
		ev.AlreadyInstalled = true
	}

	if helper.GetRecentVersion() == ev.Version {
		ev.UsedVersion = true
	}
}

func (ev ExtendedVersion) GetPromptName(showStable bool) string {
	message := ev.getCleanVersionName()

	if showStable {
		stable := "unstable"
		if ev.IsStable {
			stable = "stable"
		}

		message = fmt.Sprintf("%s (%s)", message, stable)
	}

	if ev.AlreadyInstalled && !ev.UsedVersion {
		message += " - already downloaded"
	}

	if ev.UsedVersion {
		message += " - current version"
	}

	return message
}

func (ev ExtendedVersion) getCleanVersionName() string {
	return strings.TrimPrefix(ev.Version, "go")
}

func (v Version) FilterAlreadyDownloadedVersions(evs []*ExtendedVersion) []string {
	installedVersions := make([]string, 0, len(evs))

	for _, vi := range evs {
		if vi.AlreadyInstalled {
			installedVersions = append(installedVersions, vi.Version)
		}
	}

	return installedVersions
}

func (v Version) GetVersions(forceFetchVersions bool) ([]*ExtendedVersion, error) {
	var responseVersions []api_client.VersionInfo

	if !v.fileHelpers.AreVersionsCached() || forceFetchVersions {
		err := v.clientAPI.FetchVersions(context.TODO(), &responseVersions)
		if err != nil {
			return nil, err
		}

		bytes, err := json.Marshal(responseVersions)
		if err != nil {
			return nil, err
		}

		err = v.fileHelpers.StoreVersionsResponse(bytes)
		if err != nil {
			return nil, err
		}
	} else {
		err := v.fileHelpers.GetCachedResponse(&responseVersions)
		if err != nil {
			return nil, err
		}
	}
	versions := make([]*ExtendedVersion, 0, len(responseVersions))
	for _, rv := range responseVersions {
		version := &ExtendedVersion{VersionInfo: rv}
		version.addExtras(v.fileHelpers)

		versions = append(versions, version)
	}

	return versions, nil
}

func (v Version) DeleteUnusedVersions(evs []*ExtendedVersion) (int, error) {
	versions := v.FilterAlreadyDownloadedVersions(evs)
	usedVersion := v.fileHelpers.GetRecentVersion()

	if usedVersion == "" {
		return -1, &errors.NoInstalledVersionsError{}
	}

	count := 0
	for _, version := range versions {
		if version != usedVersion {
			v.log.PrintMessage("Deleting %s.\n", version)
			if err := v.fileHelpers.DeleteDirectory(version); err != nil {
				return count, &errors.DeleteVersionError{Err: err, Version: version}
			}

			v.log.PrintMessage("%s is deleted.\n", version)

			count++
		}
	}

	return count, nil
}

func (v Version) GetLatestVersion(evs []*ExtendedVersion) int {
	for i, vi := range evs {
		if vi.IsStable {
			return i
		}
	}

	return -1
}

func (v Version) Install(ev *ExtendedVersion, os string, arch string) error {
	if ev.AlreadyInstalled {
		err := v.installer.ExistingVersion(ev.Version)
		if err != nil {
			return err
		}
	} else {
		var fileName string
		var checksum string

		for _, file := range ev.Files {
			if file.Architecture == arch && file.OS == os && file.Kind == "archive" {
				fileName = file.Filename
				checksum = file.Checksum
			}
		}

		if fileName == "" {
			return &errors.InstalledNotFoundError{OS: os, Arch: arch}
		}

		if checksum == "" {
			return &errors.ChecksumNotFoundError{OS: os, Arch: arch}
		}

		err := v.installer.NewVersion(context.TODO(), fileName, checksum, ev.Version)
		if err != nil {
			return err
		}
	}

	v.log.PrintMessage("%s version is installed!\n", ev.getCleanVersionName())

	return nil
}

func (v Version) GetPromptVersions(evs []*ExtendedVersion, showAllVersions bool) []*ExtendedVersion {
	var filteredVersions []*ExtendedVersion
	for _, version := range evs {
		if showAllVersions || (!showAllVersions && version.IsStable) {
			filteredVersions = append(filteredVersions, version)
		}
	}
	return filteredVersions
}

func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI, installer install.Installer, logger *logger.Log) Version {
	return Version{
		installer:   installer,
		clientAPI:   clientAPI,
		fileHelpers: fileHelpers,
		log:         logger,
	}
}
