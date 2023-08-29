package main

import (
	"flag"
	"runtime"

	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/version"
	"github.com/manifoldco/promptui"
)

var (
	refreshVersions = false
	installLatest   = false
)

func init() {
	files.CreateInitFiles()

	flag.BoolVar(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale")
	flag.BoolVar(&installLatest, "latest", false, "Install latest stable version")
	flag.Parse()
}

func main() {
	config := cf.GetConfig()
	versions := version.GetVersions(config.VERSIONS_URL, refreshVersions)

	var promptVersions []string
	for _, version := range versions {
		promptVersions = append(promptVersions, (&version).GetPromptName())
	}

	var selectedIndex int
	if installLatest {
		selectedIndex = version.GetLatestVersion(versions)
	} else {
		prompt := promptui.Select{
			Label: "Select go version",
			Items: promptVersions,
			Size:  10,
		}

		var errPrompt error
		selectedIndex, _, errPrompt = prompt.Run()
		if errPrompt != nil {
			panic(errPrompt)
		}
	}

	versions[selectedIndex].Install(runtime.GOOS, runtime.GOARCH, config.DOWNLOAD_VERSION_BASE_URL)
}
