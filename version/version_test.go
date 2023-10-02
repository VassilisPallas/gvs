package version_test

import (
	"fmt"
	"testing"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/version"
	"github.com/google/go-cmp/cmp"
)

func TestFilterAlreadyDownloadedVersionsReturnInstalledVersions(t *testing.T) {
	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	evs := []*version.ExtendedVersion{
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.8",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.7",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.6",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	res := versioner.FilterAlreadyDownloadedVersions(evs)

	expectedVersions := []string{
		"go1.21.0",
		"go1.20.8",
		"go1.20.7",
	}

	if !cmp.Equal(res, expectedVersions) {
		t.Errorf("Wrong array received, got=%s", cmp.Diff(res, expectedVersions))
	}
}

func TestFilterAlreadyDownloadedVersionsReturnEmptyResults(t *testing.T) {
	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	evs := []*version.ExtendedVersion{
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.8",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.7",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.6",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	res := versioner.FilterAlreadyDownloadedVersions(evs)

	expectedVersions := []string{}

	if !cmp.Equal(res, expectedVersions) {
		t.Errorf("Wrong array received, got=%s", cmp.Diff(res, expectedVersions))
	}
}

func TestGetVersionsFromRequest(t *testing.T) {
	parameters := []bool{
		false,
		true,
	}

	for _, forceFetchVersions := range parameters {
		t.Run(fmt.Sprintf("with forceFetchVersions %t", forceFetchVersions), func(t *testing.T) {
			fileHelpers := &testutils.FakeFilesHelper{}
			clientAPI := testutils.FakeGoClientAPI{}
			installer := &testutils.FakeInstaller{}
			log := logger.New(&testutils.FakeStdout{}, nil)

			versioner := version.New(fileHelpers, clientAPI, installer, log)

			versions, err := versioner.GetVersions(forceFetchVersions)

			if err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				t.Fail()
			}

			if len(versions) != 3 {
				t.Errorf("versions should be 3, instead got %d", len(versions))
				t.Fail()
			}
		})
	}
}

func TestGetVersionsRequestError(t *testing.T) {
	expectedError := fmt.Errorf("An error happened")

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{
		FetchVersionsError: expectedError,
	}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(true)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}

	if versions != nil {
		t.Error("versions should be nil")
	}
}

func TestGetVersionsFromCache(t *testing.T) {
	fileHelpers := &testutils.FakeFilesHelper{
		CachedVersion: true,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(false)

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}

	if len(versions) != 4 {
		t.Errorf("versions should be 4, instead got %d", len(versions))
	}
}

func TestGetVersionsAddCorrectExtras(t *testing.T) {
	fileHelpers := &testutils.FakeFilesHelper{
		CachedVersion:             true,
		RecentVersion:             "go1.21.0",
		AlreadyDownloadedVersions: []string{"go1.21.0", "go1.19.0"},
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(false)

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}

	if len(versions) != 4 {
		t.Errorf("versions should be 4, instead got %d", len(versions))
	}

	for _, v := range versions {
		if v.Version == "go1.21.0" {
			if !v.AlreadyInstalled {
				t.Errorf("%s should be already installed", v.Version)
			}

			if !v.UsedVersion {
				t.Errorf("%s should be the used version", v.Version)
			}
		} else if v.Version == "go1.19.0" {
			if !v.AlreadyInstalled {
				t.Errorf("%s should be already installed", v.Version)
			}

			if v.UsedVersion {
				t.Errorf("%s should not be the used version", v.Version)
			}
		} else {
			if v.AlreadyInstalled {
				t.Errorf("%s should not be already installed", v.Version)
			}

			if v.UsedVersion {
				t.Errorf("%s should not be the used version", v.Version)
			}
		}
	}
}

func TestGetVersionsFromCacheError(t *testing.T) {
	expectedError := fmt.Errorf("An error happened")

	fileHelpers := &testutils.FakeFilesHelper{
		CacheResponseError: expectedError,
		CachedVersion:      true,
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(false)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}

	if versions != nil {
		t.Error("versions should be nil")
	}
}

func TestDeleteUnusedVersionsDeleteAllUnusedVersions(t *testing.T) {
	printer := &testutils.FakeStdout{}

	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.19.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{
		RecentVersion: "go1.21.0",
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(printer, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	count, err := versioner.DeleteUnusedVersions(versions)

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}

	if count != 2 {
		t.Errorf("error should be 2, instead got %d", count)
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"Deleting go1.20.0.\n",
		"go1.20.0 is deleted.\n",
		"Deleting go1.19.0.\n",
		"go1.19.0 is deleted.\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestDeleteUnusedVersionsReturnErrorWhenNoRecentVersion(t *testing.T) {
	expectedError := fmt.Errorf("there is no any installed version")

	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.19.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{
		RecentVersion: "",
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	count, err := versioner.DeleteUnusedVersions(versions)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}

	if count != -1 {
		t.Errorf("error should be -1, instead got %d", count)
	}
}

func TestDeleteUnusedVersionsReturnErrorOnDelete(t *testing.T) {
	printer := &testutils.FakeStdout{}

	errorMessage := "some error occurred while deleting the version"

	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "bad_version",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{
		RecentVersion:        "go1.21.0",
		DeleteDirectoryError: fmt.Errorf(errorMessage),
	}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(printer, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	count, err := versioner.DeleteUnusedVersions(versions)

	expectedError := fmt.Errorf("an error occurred while deleting \"bad_version\": %q", errorMessage)
	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}

	if count != 1 {
		t.Errorf("error should be 1, instead got %d", count)
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"Deleting go1.20.0.\n",
		"go1.20.0 is deleted.\n",
		"Deleting bad_version.\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestGetLatestVersionReturnLatestStableVersion(t *testing.T) {
	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "1.22rc1",
				IsStable: false,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "1.20rc3",
				IsStable: false,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	index := versioner.GetLatestVersion(versions)

	if index != 1 {
		t.Errorf("error should be 1, instead got %d", index)
	}
}

func TestGetLatestVersionReturnLatestStableVersionNoStableVersionFound(t *testing.T) {
	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "1.22rc1",
				IsStable: false,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	index := versioner.GetLatestVersion(versions)

	if index != -1 {
		t.Errorf("error should be -1, instead got %d", index)
	}
}

func TestInstallShouldInstallExistingVersionSuccess(t *testing.T) {
	os := "darwin"
	arch := "arm64"

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: true,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files:    []api_client.FileInformation{},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if !installer.ExistingVersionCalled {
		t.Errorf("ExistingVersion should have been called")
	}

	if installer.NewVersionCalled {
		t.Errorf("NewVersion should not have been called")
	}

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}
}

func TestInstallShouldInstallExistingVersionLogs(t *testing.T) {
	printer := &testutils.FakeStdout{}

	os := "darwin"
	arch := "arm64"

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: true,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files:    []api_client.FileInformation{},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(printer, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"1.21.0 version is installed!\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestInstallExistingVersionError(t *testing.T) {
	os := "darwin"
	arch := "arm64"
	expectedError := fmt.Errorf("error while installing version")

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: true,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files:    []api_client.FileInformation{},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{
		ExistingVersionError: expectedError,
	}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionFileNameNotFound(t *testing.T) {
	os := "darwin"
	arch := "arm64"
	expectedError := fmt.Errorf("installer not found for %q %q", os, arch)

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "amd64",
					Version:      "go1.21.0",
					Checksum:     "ccd94d7a7b4f3d3e038d0ec608334c827ee8c67fc4c80a6d6037c8f5938aeb78",
					Size:         64768082,
					Kind:         "archive",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionChecksumNotFound(t *testing.T) {
	os := "darwin"
	arch := "arm64"
	expectedError := fmt.Errorf("checksum not found for %q %q", os, arch)

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "arm64",
					Version:      "go1.21.0",
					Checksum:     "",
					Size:         64768082,
					Kind:         "archive",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionArchiveNotFound(t *testing.T) {
	os := "darwin"
	arch := "arm64"
	expectedError := fmt.Errorf("installer not found for %q %q", os, arch)

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "arm64",
					Version:      "go1.21.0",
					Checksum:     "ccd94d7a7b4f3d3e038d0ec608334c827ee8c67fc4c80a6d6037c8f5938aeb78",
					Size:         64768082,
					Kind:         "source",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestInstallNewVersionSuccess(t *testing.T) {
	os := "darwin"
	arch := "arm64"

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "arm64",
					Version:      "go1.21.0",
					Checksum:     "ccd94d7a7b4f3d3e038d0ec608334c827ee8c67fc4c80a6d6037c8f5938aeb78",
					Size:         64768082,
					Kind:         "archive",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if !installer.NewVersionCalled {
		t.Errorf("NewVersionCalled should have been called")
	}

	if installer.ExistingVersionCalled {
		t.Errorf("ExistingVersion should not have been called")
	}

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}
}

func TestInstallNewVersionLogs(t *testing.T) {
	printer := &testutils.FakeStdout{}

	os := "darwin"
	arch := "arm64"

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "arm64",
					Version:      "go1.21.0",
					Checksum:     "ccd94d7a7b4f3d3e038d0ec608334c827ee8c67fc4c80a6d6037c8f5938aeb78",
					Size:         64768082,
					Kind:         "archive",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(printer, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err != nil {
		t.Errorf("error should be nil, instead got %q", err.Error())
	}

	printedMessages := printer.GetPrintMessages()
	expectedPrintedMessages := []string{
		"1.21.0 version is installed!\n",
	}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong logs received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}
}

func TestInstallNewVersionError(t *testing.T) {
	os := "darwin"
	arch := "arm64"
	expectedError := fmt.Errorf("error while installing version")

	ev := version.ExtendedVersion{
		UsedVersion:      false,
		AlreadyInstalled: false,
		VersionInfo: api_client.VersionInfo{
			Version:  "go1.21.0",
			IsStable: true,
			Files: []api_client.FileInformation{
				{
					Filename:     "go1.21.0.darwin-amd64.tar.gz",
					OS:           "darwin",
					Architecture: "arm64",
					Version:      "go1.21.0",
					Checksum:     "ccd94d7a7b4f3d3e038d0ec608334c827ee8c67fc4c80a6d6037c8f5938aeb78",
					Size:         64768082,
					Kind:         "archive",
				},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{
		NewVersionError: expectedError,
	}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	err := versioner.Install(&ev, os, arch)

	if err.Error() != expectedError.Error() {
		t.Errorf("error should be %q, instead got %q", expectedError.Error(), err.Error())
	}
}

func TestGetPromptVersionsStableOnly(t *testing.T) {
	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "1.20rc3",
				IsStable: false,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	expectedVersions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}
	filteredVersions := versioner.GetPromptVersions(versions, false)

	if !cmp.Equal(expectedVersions, filteredVersions) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(expectedVersions, filteredVersions))
	}
}

func TestGetPromptVersionsAllVersions(t *testing.T) {
	versions := []*version.ExtendedVersion{
		{
			UsedVersion:      true,
			AlreadyInstalled: true,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.21.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.20.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "1.20rc3",
				IsStable: false,
				Files:    []api_client.FileInformation{},
			},
		},
		{
			UsedVersion:      false,
			AlreadyInstalled: false,
			VersionInfo: api_client.VersionInfo{
				Version:  "go1.18.0",
				IsStable: true,
				Files:    []api_client.FileInformation{},
			},
		},
	}

	fileHelpers := &testutils.FakeFilesHelper{}
	clientAPI := testutils.FakeGoClientAPI{}
	installer := &testutils.FakeInstaller{}
	log := logger.New(&testutils.FakeStdout{}, nil)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	filteredVersions := versioner.GetPromptVersions(versions, true)

	if !cmp.Equal(versions, filteredVersions) {
		t.Errorf("Wrong object received, got=%s", cmp.Diff(versions, filteredVersions))
	}
}

func TestExtendedVersionGetPromptNameShowStable(t *testing.T) {
	parameters := []struct {
		version    version.ExtendedVersion
		showStable bool
		message    string
	}{
		{
			version: version.ExtendedVersion{
				UsedVersion:      true,
				AlreadyInstalled: true,
				VersionInfo: api_client.VersionInfo{
					Version:  "go1.21.0",
					IsStable: true,
					Files:    []api_client.FileInformation{},
				},
			},
			showStable: true,
			message:    "1.21.0 (stable) - current version",
		},
		{
			version: version.ExtendedVersion{
				UsedVersion:      false,
				AlreadyInstalled: true,
				VersionInfo: api_client.VersionInfo{
					Version:  "go1.20.0",
					IsStable: true,
					Files:    []api_client.FileInformation{},
				},
			},
			showStable: true,
			message:    "1.20.0 (stable) - already downloaded",
		},
		{
			version: version.ExtendedVersion{
				UsedVersion:      false,
				AlreadyInstalled: false,
				VersionInfo: api_client.VersionInfo{
					Version:  "go1.21rc2",
					IsStable: false,
					Files:    []api_client.FileInformation{},
				},
			},
			showStable: true,
			message:    "1.21rc2 (unstable)",
		},
		{
			version: version.ExtendedVersion{
				UsedVersion:      true,
				AlreadyInstalled: true,
				VersionInfo: api_client.VersionInfo{
					Version:  "go1.21.0",
					IsStable: true,
					Files:    []api_client.FileInformation{},
				},
			},
			showStable: false,
			message:    "1.21.0 - current version",
		},
		{
			version: version.ExtendedVersion{
				UsedVersion:      false,
				AlreadyInstalled: false,
				VersionInfo: api_client.VersionInfo{
					Version:  "go1.21rc2",
					IsStable: false,
					Files:    []api_client.FileInformation{},
				},
			},
			showStable: false,
			message:    "1.21rc2",
		},
	}

	for _, param := range parameters {
		t.Run("", func(t *testing.T) {
			res := param.version.GetPromptName(param.showStable)

			if res != param.message {
				t.Errorf("result should be %q, instead got %q", param.message, res)
				t.Fail()
			}
		})
	}
}
