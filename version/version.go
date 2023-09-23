package version

import (
	"context"
	"fmt"
	"strings"

	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/install"

	"github.com/VassilisPallas/gvs/api_client"
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
	Installer   install.Installer
	ClientAPI   api_client.GoClientAPI
	FileHelpers files.FileHelpers

	Versioner
}

type ExtendedVersion struct {
	UsedVersion      bool
	AlreadyInstalled bool

	api_client.VersionInfo
}

// TODO: test this
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

	if !v.FileHelpers.AreVersionsCached() || forceFetchVersions {
		err := v.ClientAPI.FetchVersions(context.TODO(), &responseVersions)
		if err != nil {
			return nil, err
		}
	} else {
		err := v.FileHelpers.GetCachedResponse(&responseVersions)
		if err != nil {
			return nil, err
		}
	}

	versions := make([]*ExtendedVersion, 0, len(responseVersions))
	for _, rv := range responseVersions {
		version := &ExtendedVersion{VersionInfo: rv}
		version.addExtras(v.FileHelpers)

		versions = append(versions, version)
	}

	return versions, nil
}

func (v Version) DeleteUnusedVersions(evs []*ExtendedVersion) (int, error) {
	versions := v.FilterAlreadyDownloadedVersions(evs)
	usedVersion := v.FileHelpers.GetRecentVersion()

	if usedVersion == "" {
		return -1, fmt.Errorf("there is no any installed version")
	}

	count := 0
	for _, version := range versions {
		if version != usedVersion {
			fmt.Printf("Deleting %s \n", version)
			if err := v.FileHelpers.DeleteDirectory(version); err != nil {
				return count, fmt.Errorf("an error occurred while deleting %s: %s", version, err.Error())
			}

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
		// TODO: validate this got called
		err := v.Installer.ExistingVersion(ev.Version)
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
			return fmt.Errorf("installer not found for %s-%s", os, arch)
		}

		if checksum == "" {
			return fmt.Errorf("checksum not found for %s-%s", os, arch)
		}

		// TODO: validate this got called
		err := v.Installer.NewVersion(context.TODO(), fileName, checksum, ev.Version)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s version is installed!\n", ev.getCleanVersionName())

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

func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI, installer install.Installer) Versioner {
	return Version{
		Installer:   installer,
		ClientAPI:   clientAPI,
		FileHelpers: fileHelpers,
	}
}
