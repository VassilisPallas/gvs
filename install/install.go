package install

import (
	"fmt"
	"io"
	"net/http"

	"github.com/VassilisPallas/gvs/files"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func downloadVersionFile(url string, downloadChannel chan<- io.ReadCloser) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	response, err := Client.Do(request)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != 200 {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	downloadChannel <- response.Body
}

func compareChecksums(checksum string) error {
	sha256sum, err := files.GetTarChecksum()
	if err != nil {
		files.RemoveTarFile()
		panic(err)
	}

	hash := fmt.Sprintf("%x", sha256sum)

	if hash != checksum {
		return fmt.Errorf("checksums do not match.\nExpected: %s\nGot: %s", checksum, hash)
	}

	return nil
}

func NewVersion(baseURL string, fileName string, checksum string, goVersionName string) {
	downloadChannel := make(chan io.ReadCloser)

	fmt.Println("Downloading...")
	url := fmt.Sprintf("%s/%s", baseURL, fileName)
	go downloadVersionFile(url, downloadChannel)

	content := <-downloadChannel
	defer content.Close()

	err := files.CreateTarFile(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("Compare Checksums...")
	if err := compareChecksums(checksum); err != nil {
		panic(err)
	}

	fmt.Println("Unzipping...")
	if err := files.UnzipTarFile(); err != nil {
		panic(err)
	}

	files.RenameGoDirectory(goVersionName)
	files.RemoveTarFile()

	fmt.Println("Installing version...")
	if err := files.CreateExecutableSymlink(goVersionName); err != nil {
		panic(err)
	}

	files.UpdateRecentVersion(goVersionName)
}

func ExistingVersion(goVersionName string) {
	fmt.Println("Installing version...")
	files.CreateExecutableSymlink(goVersionName)
	files.UpdateRecentVersion(goVersionName)
}
