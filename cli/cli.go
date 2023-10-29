package cli

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/VassilisPallas/gvs/logger"
	"github.com/VassilisPallas/gvs/version"
)

type CLI struct {
	versions  []*version.ExtendedVersion
	versioner version.Versioner
	log       logger.Logger
}

func (cli CLI) Install(selectedVersion *version.ExtendedVersion) error {
	cli.log.Info("selected %s version\n", selectedVersion.Version)
	return cli.versioner.Install(selectedVersion, runtime.GOOS, runtime.GOARCH)
}

func (cli CLI) InstallVersion(goVersion string) error {
	semver := &version.Semver{}
	err := version.ParseSemver(goVersion, semver)
	if err != nil {
		return err
	}

	selectedVersion := cli.versioner.FindVersionBasedOnSemverName(cli.versions, semver)
	if selectedVersion == nil {
		return fmt.Errorf("%s is not a valid version", semver.GetVersion())
	}

	return cli.Install(selectedVersion)
}

func (cli CLI) InstallLatestVersion() error {
	selectedIndex := cli.versioner.GetLatestVersion(cli.versions)
	if selectedIndex == -1 {
		return errors.New("latest version not found")
	}

	selectedVersion := cli.versions[selectedIndex]

	return cli.Install(selectedVersion)
}

func (cli CLI) DeleteUnusedVersions() error {
	deleted_count, err := cli.versioner.DeleteUnusedVersions(cli.versions)
	if err != nil {
		return err
	}

	if deleted_count > 0 {
		cli.log.PrintMessage("All the unused version are deleted!")
	} else {
		cli.log.PrintMessage("Nothing to delete")
	}

	return nil
}

func New(versions []*version.ExtendedVersion, versioner version.Versioner, log logger.Logger) CLI {
	return CLI{versions: versions, versioner: versioner, log: log}
}
