package main

import (
	"flag"
	"fmt"
	"runtime"

	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/version"
	"github.com/manifoldco/promptui"
)

var (
	refreshVersions = false
	installLatest   = false
	deleteUnused    = false
	showAllVersions = false
)

func init() {
	files.CreateInitFiles()

	flag.BoolVar(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale")
	flag.BoolVar(&installLatest, "latest", false, "Install latest stable version")
	flag.BoolVar(&deleteUnused, "delete-unused", false, "Delete all unused versions")
	flag.BoolVar(&showAllVersions, "all", false, "Show both stable and unstable versions")
	flag.Parse()
}

func main() {
	config := cf.GetConfig()

	versions := version.GetVersions(config.GO_BASE_URL, refreshVersions, showAllVersions || deleteUnused)

	switch {
	case deleteUnused:
		deleted_count := version.DeleteUnusedVersions(version.FilterAlreadyDownloadedVersions(versions))
		if deleted_count > 0 {
			fmt.Println("All the unused version are deleted!")
		} else {
			fmt.Println("Nothing to delete")
		}
	case installLatest:
		selectedIndex := version.GetLatestVersion(versions)
		versions[selectedIndex].Install(runtime.GOOS, runtime.GOARCH, config.GO_BASE_URL)
	default:
		var promptVersions []string
		for _, version := range versions {
			promptVersions = append(promptVersions, version.GetPromptName(showAllVersions))
		}

		prompt := promptui.Select{
			Label: "Select go version",
			Items: promptVersions,
			Size:  10,
		}

		selectedIndex, _, errPrompt := prompt.Run()
		if errPrompt != nil {
			panic(errPrompt)
		}
		versions[selectedIndex].Install(runtime.GOOS, runtime.GOARCH, config.GO_BASE_URL)
	}
}
