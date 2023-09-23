package files_test

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/VassilisPallas/gvs/files"
)

func TestGetHomeDirectory(t *testing.T) {
	fileUtils := files.NewUtils()

	res := fileUtils.GetHomeDirectory()

	if res == "" {
		t.Error("GetHomeDirectory should return an non-empty string")
	}
}

func TestGetAppDir(t *testing.T) {
	fileUtils := files.NewUtils()
	user, _ := user.Current()

	expected := fmt.Sprintf("%s/.gvs", user.HomeDir)
	res := fileUtils.GetAppDir()

	if res != expected {
		t.Errorf("GetAppDir should be '%s', got '%s'", expected, res)
	}
}

func TestGetVersionsDir(t *testing.T) {
	fileUtils := files.NewUtils()
	user, _ := user.Current()

	expected := fmt.Sprintf("%s/.gvs/.go.versions", user.HomeDir)
	res := fileUtils.GetVersionsDir()

	if res != expected {
		t.Errorf("GetVersionsDir should be '%s', got '%s'", expected, res)
	}
}

func TestGetTarDir(t *testing.T) {
	fileUtils := files.NewUtils()
	user, _ := user.Current()

	expected := fmt.Sprintf("%s/.gvs/.go.versions/downloaded.tar.gz", user.HomeDir)
	res := fileUtils.GetTarFile()

	if res != expected {
		t.Errorf("GetTarFile should be '%s', got '%s'", expected, res)
	}
}

func TestGetBinDir(t *testing.T) {
	fileUtils := files.NewUtils()
	user, _ := user.Current()

	expected := fmt.Sprintf("%s/bin", user.HomeDir)
	res := fileUtils.GetBinDir()

	if res != expected {
		t.Errorf("GetBinDir should be '%s', got '%s'", expected, res)
	}
}

func TestGetCurrentVersionFile(t *testing.T) {
	fileUtils := files.NewUtils()
	user, _ := user.Current()

	expected := fmt.Sprintf("%s/.gvs/.go.versions/CURRENT", user.HomeDir)
	res := fileUtils.GetCurrentVersionFile()

	if res != expected {
		t.Errorf("GetCurrentVersionFile should be '%s', got '%s'", expected, res)
	}
}

func TestGetVersionResponseFile(t *testing.T) {
	fileUtils := files.NewUtils()

	expected := "goVersions.json"
	res := fileUtils.GetVersionResponseFile()

	if res != expected {
		t.Errorf("GetVersionResponseFile should be '%s', got '%s'", expected, res)
	}
}
