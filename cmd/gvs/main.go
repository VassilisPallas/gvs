package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	terminalColors "github.com/fatih/color"

	"github.com/VassilisPallas/gvs/api_client"
	"github.com/VassilisPallas/gvs/clock"
	cf "github.com/VassilisPallas/gvs/config"
	"github.com/VassilisPallas/gvs/files"
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
)

func getBold() *terminalColors.Color {
	return terminalColors.New().Add(terminalColors.Bold)
}

func parseFlags() {
	flag.BoolVar(&showAllVersions, "show-all", false, "Show both stable and unstable versions.")
	flag.BoolVar(&installLatest, "install-latest", false, "Install latest stable version.")
	flag.BoolVar(&deleteUnused, "delete-unused", false, "Delete all unused versions that were installed before.")
	flag.BoolVar(&refreshVersions, "refresh-versions", false, "Fetch again go versions in case the cached ones are stale.")

	flag.Usage = func() {
		bold := getBold()
		gvsMessage := bold.Sprint("gvs")

		flagSet := flag.CommandLine

		fmt.Println()
		bold.Println("NAME")
		fmt.Printf("  gvs\n\n")

		bold.Println("DESCRIPTION")
		fmt.Printf("  the %s CLI is a command line tool to manage multiple active Go versions.\n\n", gvsMessage)

		bold.Println("SYNOPSIS")
		fmt.Printf("  gvs\n   [--show-all]\n   [--install-latest]\n   [--delete-unused]\n   [--refresh-versions]\n\n")

		bold.Println("FLAGS")
		flags := []string{"show-all", "install-latest", "delete-unused", "refresh-versions"}
		for _, name := range flags {
			flag := flagSet.Lookup(name)
			fmt.Printf("  --%s\n\t%s\n", flag.Name, flag.Usage)
		}
		fmt.Println()

		fmt.Printf("Before start using the %s CLI, make sure to delete all the existing go versions\n", gvsMessage)
		fmt.Printf("and append to your profile file the export: %q.\n", "export PATH=$HOME/bin:$PATH")
		fmt.Printf("The profile file could be one of: (%s)\n", "~/.bash_profile, ~/.zshrc, ~/.profile, or ~/.bashrc")
	}

	flag.Parse()
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
