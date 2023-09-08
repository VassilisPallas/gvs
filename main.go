package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/VassilisPallas/gvs/api_client"
	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/version"
	"github.com/manifoldco/promptui"
)

var (
	refreshVersions = false
	installLatest   = false
	deleteUnused    = false
	showAllVersions = false

	filesUtils = files.New()
)

func parseFlags() {
	flag.BoolVar(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale")
	flag.BoolVar(&installLatest, "latest", false, "Install latest stable version")
	flag.BoolVar(&deleteUnused, "delete-unused", false, "Delete all unused versions")
	flag.BoolVar(&showAllVersions, "all", false, "Show both stable and unstable versions")
	flag.Parse()
}

func init() {
	if err := filesUtils.CreateInitFiles(); err != nil {
		// TODO: pass error as log to file
		fmt.Println(err)
		os.Exit(1)
	}

	parseFlags()
}

func main() {
	config := cf.GetConfig()
	clientAPI := api_client.New(config)
	installer := install.New(filesUtils, clientAPI)

	versioner := version.New(filesUtils, clientAPI, installer)

	versions, err := versioner.GetVersions(refreshVersions)
	if err != nil {
		// TODO: pass error as log to file
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case deleteUnused:
		deleted_count, err := versioner.DeleteUnusedVersions(versions)
		if err != nil {
			// TODO: pass error as log to file
			fmt.Println(err)
			os.Exit(1)
		}

		if deleted_count > 0 {
			fmt.Println("All the unused version are deleted!")
		} else {
			fmt.Println("Nothing to delete")
		}
	case installLatest:
		selectedIndex := versioner.GetLatestVersion(versions)
		err := versioner.Install(versions[selectedIndex], runtime.GOOS, runtime.GOARCH, config.GO_BASE_URL)
		if err != nil {
			// TODO: pass error as log to file
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		promptVersions := versioner.GetPromptVersions(versions, showAllVersions)

		var versionNames []string

		for _, pv := range promptVersions {
			versionNames = append(versionNames, pv.GetPromptName(showAllVersions))
		}

		prompt := promptui.Select{
			Label: "Select go version",
			Items: versionNames,
			Size:  10,
		}

		selectedIndex, _, errPrompt := prompt.Run()
		if errPrompt != nil {
			// TODO: pass error as log to file
			fmt.Println(errPrompt)
			os.Exit(1)
		}
		err := versioner.Install(promptVersions[selectedIndex], runtime.GOOS, runtime.GOARCH, config.GO_BASE_URL)
		if err != nil {
			// TODO: pass error as log to file
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
