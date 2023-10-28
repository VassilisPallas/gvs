// Package version provides an interface to make handle
// the CLI logic for the versions.
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

// Versioner is the interface that wraps the basic methods for handling
// the CLI logic for the versions.
type Versioner interface {
	// GetVersions returns back a slice of versions.
	// FetchVersions must return a slice with the versions a non-null error.
	GetVersions(forceFetchVersions bool) ([]*ExtendedVersion, error)

	// DeleteUnusedVersions deletes all the unused versions.
	// The input should contain the versions that the method will iterate to find
	// and delete the unused versions.
	// DeleteUnusedVersions must return the count of the deleted versions and a non-null
	// error if the versions are deleted.
	DeleteUnusedVersions(evs []*ExtendedVersion) (int, error)

	// GetLatestVersion returns the latest stable version.
	// The input should contain the versions that the method will iterate to find and the version.
	// GetLatestVersion must return the the index of the found version, or -1 if not found.
	GetLatestVersion(evs []*ExtendedVersion) int

	// Install installs the given version for the OS and the architecture type.
	// Install must return a non-null error if the install was successful.
	Install(ev *ExtendedVersion, os string, arch string) error

	// GetPromptVersions returns a filtered list of versions based on if the version is stable or not.
	// GetPromptVersions must return a slice of *ExtendedVersion.
	GetPromptVersions(evs []*ExtendedVersion, showAllVersions bool) []*ExtendedVersion

	// FilterAlreadyDownloadedVersions returns the versions that are already installed.
	// FilterAlreadyDownloadedVersions must return a slice of string.
	FilterAlreadyDownloadedVersions(evs []*ExtendedVersion) []string

	// FindVersionBasedOnSemverName returns the version that is described in the semver.
	// In case one of minor or patch version number is missing from the semver,
	// FindVersionBasedOnSemverName should return the latest value of them.
	// If the version is not found, FindVersionBasedOnSemverName should return nil.
	FindVersionBasedOnSemverName(evs []*ExtendedVersion, version *Semver) *ExtendedVersion
}

// Version is the struct that implements the Versioner interface.
type Version struct {
	// installer is used to handle the install of the version (either for a new version, or an already downloaded one).
	installer install.Installer

	// clientAPI is used to download the selected version (if it's a new version)
	clientAPI api_client.GoClientAPI

	// fileHelpers is used to access and write files.
	fileHelpers files.FileHelpers

	// log is the custom Logger
	log *logger.Log
}

// ExtendedVersion embeds VersionInfo from the api_client.
// It adds more attributes that can be used to identify if a version is already installed and/or currently used.
type ExtendedVersion struct {
	// UsedVersion indicates if the version is currently used.
	UsedVersion bool

	// AlreadyInstalled indicates if the version is already installed.
	AlreadyInstalled bool

	api_client.VersionInfo
}

// addExtras updates the attributes based on whether the version is already installed and/or currently used.
//
// addExtras used the FileHelpers to access the requires files.
func (ev *ExtendedVersion) addExtras(helper files.FileHelpers) {
	if helper.DirectoryExists(ev.Version) {
		ev.AlreadyInstalled = true
	}

	if helper.GetRecentVersion() == ev.Version {
		ev.UsedVersion = true
	}
}

// GetPromptName returns the version name as it will be rendered on the dropdown prompt.
//
// If the `showStable` is set to true, it will also include if the version is stable or not.
// Examples:
//   - 1.21.3 (stable) - current version
//   - 1.21.0 (stable) - already downloaded
//   - 1.21rc4 (unstable)
//
// If the `showStable` is set to false, the text in the paragraphs will be omitted.
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

// getCleanVersionName removed the `go` prefix from the version name.
//
// Example:
//
// `go1.21.3` will be returns as `1.21.3`.
func (ev ExtendedVersion) getCleanVersionName() string {
	return strings.TrimPrefix(ev.Version, "go")
}

// FilterAlreadyDownloadedVersions returns the version names that are already installed.
func (v Version) FilterAlreadyDownloadedVersions(evs []*ExtendedVersion) []string {
	installedVersions := make([]string, 0, len(evs))

	for _, vi := range evs {
		if vi.AlreadyInstalled {
			installedVersions = append(installedVersions, vi.Version)
		}
	}

	return installedVersions
}

// GetVersions returns a slice with the available versions to be installed.
//
// GetVersions will first check if the response is cached. If it is, it will read it from the file.
// Otherwise, if the make the request to the API, parses the result and stores it into the cache.
//
// There is also the option to force the fetch from the API, which will re-download the versions.
//
// Finally, for each available version, GetVersions includes the "extra" attributes in each of the ExtendedVersion
// types that indicates if a version is already installed and/or currently used.
func (v Version) GetVersions(forceFetchVersions bool) ([]*ExtendedVersion, error) {
	var responseVersions []api_client.VersionInfo

	if !v.fileHelpers.AreVersionsCached() || forceFetchVersions {
		err := v.clientAPI.FetchVersions(context.Background(), &responseVersions)
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

// DeleteUnusedVersions deletes all the unused versions.
// If there is no any unused version, DeleteUnusedVersions will return -1 as the count and an error
// of the type *NoInstalledVersionsError.
//
// If an error occurs while deleting a version, DeleteUnusedVersions will return the count of the
// version that have already been deleted and an error of type *DeleteVersionError.
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

// GetLatestVersion returns the latest stable version.
//
// GetLatestVersion returns the the index of the found version, or -1 if not found.
func (v Version) GetLatestVersion(evs []*ExtendedVersion) int {
	for i, vi := range evs {
		if vi.IsStable {
			return i
		}
	}

	return -1
}

// Install installs the given version for the OS and the architecture type.
//
// If for a new version (which has not been downloaded already in the past):
//
//		if the archive file is not found for the OS and the architecture type, then an error of
//		the type *InstallerNotFoundError is returned.
//
//		If checksum is not found for the OS and the architecture type, then an error of
//	    the type *ChecksumNotFoundError is returned.
//
// Otherwise, any other error that might occur during the install of an existing or a new version
// will be returned back.
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
			return &errors.InstallerNotFoundError{OS: os, Arch: arch}
		}

		if checksum == "" {
			return &errors.ChecksumNotFoundError{OS: os, Arch: arch}
		}

		err := v.installer.NewVersion(context.Background(), fileName, checksum, ev.Version)
		if err != nil {
			return err
		}
	}

	v.log.PrintMessage("%s version is installed!\n", ev.getCleanVersionName())

	return nil
}

// GetPromptVersions returns a filtered list of versions based on if the version is stable or not.
func (v Version) GetPromptVersions(evs []*ExtendedVersion, showAllVersions bool) []*ExtendedVersion {
	var filteredVersions []*ExtendedVersion
	for _, version := range evs {
		if showAllVersions || (!showAllVersions && version.IsStable) {
			filteredVersions = append(filteredVersions, version)
		}
	}
	return filteredVersions
}

// FindVersionBasedOnSemverName returns the version that is described in the semver.
// It compares the stringified semver version as a prefix for each one of the versions.
// FindVersionBasedOnSemverName returns back the first occurrence of that version, which ensures the
// latest version will be returned when minor or patch versions are not assigned to the semver.
// If the version is not found, FindVersionBasedOnSemverName return back nil.
func (v Version) FindVersionBasedOnSemverName(evs []*ExtendedVersion, version *Semver) *ExtendedVersion {
	expectedVersion := version.GetVersion()

	for _, ev := range evs {
		if strings.HasPrefix(ev.getCleanVersionName(), expectedVersion) {
			return ev
		}
	}

	return nil
}

// New returns a Version instance that implements the Versioner interface.
// Each call to New returns a distinct Version instance even if the parameters are identical.
func New(fileHelpers files.FileHelpers, clientAPI api_client.GoClientAPI, installer install.Installer, logger *logger.Log) Version {
	return Version{
		installer:   installer,
		clientAPI:   clientAPI,
		fileHelpers: fileHelpers,
		log:         logger,
	}
}
