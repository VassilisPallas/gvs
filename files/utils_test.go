package files

import (
	"testing"

	"github.com/VassilisPallas/gvs/internal/testutils"
)

func TestGetAppDir(t *testing.T) {
	fs := testutils.FakeFileSystem{HomeDir: "/Users/someone"}

	res := getAppDir(fs)
	expectedResult := "/Users/someone/.gvs"
	if res != expectedResult {
		t.Errorf("application directory should be %q, instead got %q", expectedResult, res)
	}
}

func TestGetVersionDir(t *testing.T) {
	fs := testutils.FakeFileSystem{HomeDir: "/Users/someone"}

	res := getVersionsDir(fs)
	expectedResult := "/Users/someone/.gvs/.go.versions"
	if res != expectedResult {
		t.Errorf("application directory should be %q, instead got %q", expectedResult, res)
	}
}

func TestGetTarFile(t *testing.T) {
	fs := testutils.FakeFileSystem{HomeDir: "/Users/someone"}

	res := getTarFile(fs)
	expectedResult := "/Users/someone/.gvs/.go.versions/downloaded.tar.gz"
	if res != expectedResult {
		t.Errorf("application directory should be %q, instead got %q", expectedResult, res)
	}
}

func TestGetBinDir(t *testing.T) {
	fs := testutils.FakeFileSystem{HomeDir: "/Users/someone"}

	res := getBinDir(fs)
	expectedResult := "/Users/someone/bin"
	if res != expectedResult {
		t.Errorf("application directory should be %q, instead got %q", expectedResult, res)
	}
}

func TestGetCurrectVersionFile(t *testing.T) {
	fs := testutils.FakeFileSystem{HomeDir: "/Users/someone"}

	res := getCurrentVersionFile(fs)
	expectedResult := "/Users/someone/.gvs/.go.versions/CURRENT"
	if res != expectedResult {
		t.Errorf("application directory should be %q, instead got %q", expectedResult, res)
	}
}
