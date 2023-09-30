package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/VassilisPallas/gvs/api_client"
	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/version"
	"github.com/manifoldco/promptui"
)

var (
	refreshVersions = false
	installLatest   = false
	deleteUnused    = false
	showAllVersions = false

	filesUtils             = files.NewUtils()
	log        *logger.Log = nil
)

func parseFlags() {
	flag.BoolVar(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale")
	flag.BoolVar(&installLatest, "latest", false, "Install latest stable version")
	flag.BoolVar(&deleteUnused, "delete-unused", false, "Delete all unused versions")
	flag.BoolVar(&showAllVersions, "all", false, "Show both stable and unstable versions")
	flag.Parse()
}

func init() {
	logFile, err := filesUtils.CreateLogFile()
	log = logger.New(os.Stdout, logFile)
	if err != nil {
		log.PrintError(err.Error())
		os.Exit(1)
		return
	}

	if err := filesUtils.CreateInitFiles(); err != nil {
		log.PrintError(err.Error())
		os.Exit(1)
		return
	}

	parseFlags()
}

func main() {
	// close log file after the execution
	defer log.Close()

	config := cf.GetConfig()
	clientAPI := api_client.New(config)
	fileHelpers := files.New(filesUtils)
	installer := install.New(fileHelpers, clientAPI, log)

	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(refreshVersions)
	if err != nil {
		log.PrintError(err.Error())
		os.Exit(1)
		return
	}

	switch {
	case deleteUnused:
		log.Info("deleteUnused option selected")

		deleted_count, err := versioner.DeleteUnusedVersions(versions)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}

		if deleted_count > 0 {
			log.PrintMessage("All the unused version are deleted!")
		} else {
			log.PrintMessage("Nothing to delete")
		}
	case installLatest:
		log.Info("installLatest option selected")

		selectedIndex := versioner.GetLatestVersion(versions)
		selectedVersion := versions[selectedIndex]

		log.Info("selected %s version", selectedVersion.Version)
		err := versioner.Install(selectedVersion, runtime.GOOS, runtime.GOARCH)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	default:
		log.Info("install version option selected\n")

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
			log.PrintError(errPrompt.Error())
			os.Exit(1)
			return
		}

		selectedVersion := promptVersions[selectedIndex]

		log.Info("selected %s version\n", selectedVersion.Version)
		err := versioner.Install(selectedVersion, runtime.GOOS, runtime.GOARCH)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	}
}
