package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/install"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FileInformation struct {
	Filename     string `json:"filename"`
	OS           string `json:"os"`
	Architecture string `json:"arch"`
	Version      string `json:"version"`
	Checksum     string `json:"sha256"`
	Size         int    `json:"size"`
	Kind         string `json:"kind"`
}

type VersionInfo struct {
	Version          string            `json:"version"`
	IsStable         bool              `json:"stable"`
	Files            []FileInformation `json:"files"`
	UsedVersion      bool
	AlreadyInstalled bool
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func (vi *VersionInfo) AddExtras() {

	if files.DirectoryExists(vi.Version) {
		vi.AlreadyInstalled = true
	}

	if files.GetRecentVersion() == vi.Version {
		vi.UsedVersion = true
	}
}

func (vi VersionInfo) GetPromptName(showStable bool) string {
	message := getCleanVersionName(vi.Version)

	if showStable {
		stable := "unstable"
		if vi.IsStable {
			stable = "stable"
		}

		message = fmt.Sprintf("%s (%s)", message, stable)
	}

	if vi.AlreadyInstalled && !vi.UsedVersion {
		message += " - already downloaded"
	}

	if vi.UsedVersion {
		message += " - current version"
	}

	return message
}

func (vi VersionInfo) Install(os string, arch string, downloadURL string) {
	if vi.AlreadyInstalled {
		install.ExistingVersion(vi.Version)
	} else {
		var fileName string
		var checksum string

		for _, file := range vi.Files {
			if file.Architecture == arch && file.OS == os && file.Kind == "archive" {
				fileName = file.Filename
				checksum = file.Checksum
			}
		}

		if fileName == "" {
			panic(fmt.Errorf("installer not found for %s-%s.", os, arch))
		}

		install.NewVersion(downloadURL, fileName, checksum, vi.Version)
	}

	fmt.Printf("%s version is installed!\n", getCleanVersionName(vi.Version))
}

func getCleanVersionName(version string) string {
	return strings.TrimPrefix(version, "go")
}

func fetchVersions(baseURL string, versionsChannel chan<- []byte) {
	url := fmt.Sprintf("%s/?mode=json&include=all", baseURL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	response, err := Client.Do(request)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != 200 {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	versionsChannel <- body
}

func filterVersions(versions []*VersionInfo, showAllVersions bool) []*VersionInfo {
	var filteredVersions []*VersionInfo
	for _, version := range versions {
		if showAllVersions || (!showAllVersions && version.IsStable) {
			filteredVersions = append(filteredVersions, version)
		}
	}
	return filteredVersions
}

func GetVersions(url string, forceFetchVersions bool, showAllVersions bool) []*VersionInfo {
	var byte_versions []byte

	if files.AreVersionsCached() || forceFetchVersions {
		versionsChannel := make(chan []byte)

		go fetchVersions(url, versionsChannel)

		byte_versions = <-versionsChannel
		close(versionsChannel)

		if err := files.StoreVersionsResponse(byte_versions); err != nil {
			panic(err)
		}
	} else {
		var err error
		byte_versions, err = files.GetVersionsResponse()
		if err != nil {
			panic(err)
		}
	}

	var versions []*VersionInfo
	if err := json.Unmarshal(byte_versions, &versions); err != nil {
		panic(err)
	}

	for _, vi := range versions {
		vi.AddExtras()
	}

	return filterVersions(versions, showAllVersions)
}

func GetLatestVersion(vis []*VersionInfo) int {
	for i, vi := range vis {
		if vi.IsStable {
			return i
		}
	}

	return -1
}

func FilterAlreadyDownloadedVersions(vis []*VersionInfo) []string {
	installedVersions := make([]string, 0, len(vis))

	for _, vi := range vis {
		if vi.AlreadyInstalled {
			installedVersions = append(installedVersions, vi.Version)
		}
	}

	return installedVersions
}

func DeleteUnusedVersions(versions []string) int {
	usedVersion := files.GetRecentVersion()

	if usedVersion == "" {
		panic(fmt.Errorf("There is no any installed version"))
	}

	count := 0
	for _, version := range versions {
		if version != usedVersion {
			count++
			fmt.Printf("Deleting %s \n", version)
			files.DeleteDirectory(version)
		}
	}

	return count
}
