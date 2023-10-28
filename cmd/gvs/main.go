package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/clock"
	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
	"github.com/VassilisPallas/gvs/flags"
	"github.com/VassilisPallas/gvs/install"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/pkg/unzip"
	"github.com/VassilisPallas/gvs/version"
	"github.com/manifoldco/promptui"
)

var (
	refreshVersions = false
	installLatest   = false
	deleteUnused    = false
	showAllVersions = false
	specificVersion = ""
)

func parseFlags() {
	set := flags.FlagSet{}
	set.FlagBool(&showAllVersions, "show-all", false, "Show both stable and unstable versions.")
	set.FlagBool(&installLatest, "install-latest", false, "Install latest stable version.")
	set.FlagBool(&deleteUnused, "delete-unused", false, "Delete all unused versions that were installed before.")
	set.FlagBool(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale.")
	set.FlagStr(&specificVersion, "install-version", "", "Pass the version you want to install instead of selecting from the dropdown. If you do not specify the minor or the patch version, the latest one will be selected.")

	set.Parse()
}

func main() {
	config := cf.GetConfig()
	log := logger.New(os.Stdout, nil)

	fs := files.FileSystem{}
	unzipper := unzip.Unzip{FileSystem: fs}
	realClock := clock.RealClock{}
	fileHelpers := files.New(fs, realClock, unzipper, log)

	logFile, err := fileHelpers.CreateInitFiles()

	if err != nil {
		// this will be print only on the terminal since
		// the logger output is nil
		log.PrintError(err.Error())
		os.Exit(1)
		return
	}

	log.SetLogWriter(logFile)
	defer log.Close() // close log file after the execution

	parseFlags()

	httpClient := &http.Client{
		Timeout: time.Duration(config.REQUEST_TIMEOUT) * time.Second,
	}

	clientAPI := api_client.New(httpClient, config.GO_BASE_URL)
	installer := install.New(fileHelpers, clientAPI, log)
	versioner := version.New(fileHelpers, clientAPI, installer, log)

	versions, err := versioner.GetVersions(refreshVersions)
	if err != nil {
		log.PrintError(err.Error())
		os.Exit(1)
		return
	}

	switch {
	case specificVersion != "":
		semver := &version.Semver{}
		err := version.ParseSemver(specificVersion, semver)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}

		selectedVersion := versioner.FindVersionBasedOnSemverName(versions, semver)
		if selectedVersion == nil {
			log.PrintError("%s is not a valid version.", semver.GetVersion())
			os.Exit(1)
			return
		}

		err = versioner.Install(selectedVersion, runtime.GOOS, runtime.GOARCH)

		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
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
		if selectedIndex == -1 {
			log.PrintError("latest version not found")
			os.Exit(1)
			return
		}

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
