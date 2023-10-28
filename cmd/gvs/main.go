package main

import (
	"net/http"
	"os"
	"time"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/cli"
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
	fromModFile     = false
	specificVersion = ""
)

func parseFlags() {
	set := flags.FlagSet{}
	set.FlagBool(&showAllVersions, "show-all", 'a', false, "Show both stable and unstable versions.")
	set.FlagBool(&installLatest, "install-latest", 'l', false, "Install latest stable version.")
	set.FlagBool(&deleteUnused, "delete-unused", 'd', false, "Delete all unused versions that were installed before.")
	set.FlagBool(&refreshVersions, "refresh-versions", 'r', false, "Fetch again go versions in case the cached ones are stale.")
	set.FlagStr(&specificVersion, "install-version", 'v', "", "Pass the version you want to install instead of selecting from the dropdown. If you do not specify the minor or the patch version, the latest one will be selected.")
	set.FlagBool(&fromModFile, "from-mod", 'm', false, "Install the version that will be found on the go.mod file. The go.mod file should be on the same path you run gvs. If the version in the go.mod file do not specify the minor or the patch version, the latest one will be selected.")

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

	cli := cli.New(versions, versioner, log)

	switch {
	case fromModFile:
		log.Info("install version from go.mod file option selected")

		version, err := fileHelpers.ReadVersionFromMod()
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}

		err = cli.InstallVersion(version)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	case specificVersion != "":
		log.Info("install specific version option selected")

		err := cli.InstallVersion(specificVersion)
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	case installLatest:
		log.Info("install latest option selected")

		err := cli.InstallLatestVersion()
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

		err = cli.Install(promptVersions[selectedIndex])
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	}

	// execute this command at the end, after any other selected flag
	// to delete the unused version if needed.
	if deleteUnused {
		log.Info("deleteUnused option selected")

		err := cli.DeleteUnusedVersions()
		if err != nil {
			log.PrintError(err.Error())
			os.Exit(1)
			return
		}
	}
}
