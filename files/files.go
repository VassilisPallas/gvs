package files

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	appDir              string = ".gvs"
	versionResponseFile string = "goVersions.json"
	tarFileName         string = "downloaded.tar.gz"
	goVersionsDir       string = ".go.versions"
	binDir              string = "bin"
)

func getBaseDir() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.HomeDir
}

func getAppDir() string {
	return fmt.Sprintf("%s/%s", getBaseDir(), appDir)
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func getVersionsDir() string {
	return fmt.Sprintf("%s/%s", getAppDir(), goVersionsDir)
}

func getTarFile() string {
	return fmt.Sprintf("%s/%s", getVersionsDir(), tarFileName)
}

func getBinDir() string {
	return fmt.Sprintf("%s/%s", getBaseDir(), binDir)
}

func getCurrentVersionFile() string {
	return fmt.Sprintf("/%s/CURRENT", getVersionsDir())
}

func CreateInitFiles() {
	createDirIfNotExist(getAppDir())
	createDirIfNotExist(getVersionsDir())
	createDirIfNotExist(getBinDir())
}

func StoreVersionsResponse(body []byte) error {
	return os.WriteFile(fmt.Sprintf("%s/%s", getAppDir(), versionResponseFile), body, 0644)
}

func GetVersionsResponse() ([]byte, error) {
	return os.ReadFile(fmt.Sprintf("%s/%s", getAppDir(), versionResponseFile))
}

func AreVersionsCached() bool {
	file := fmt.Sprintf("%s/%s", getAppDir(), versionResponseFile)

	if info, err := os.Stat(file); err == nil {
		currentTime := time.Now()
		// even if the versions are cached, we return false if the
		// the date the response was stored is more than a week
		return currentTime.Sub(info.ModTime()).Hours() > 24*7
	}

	return true
}

func CreateTarFile(content io.ReadCloser) error {
	file, err := os.Create(getTarFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func GetTarChecksum() ([]byte, error) {
	hasher := sha256.New()
	path := getTarFile()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func UnzipTarFile() error {
	target := getVersionsDir()

	reader, err := os.Open(getTarFile())
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(path, 0770); err != nil {
				return err
			}

			continue
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
				return err
			}
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 777)
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}

		if err := uncompressedStream.Close(); err != nil {
			panic(err)
		}
	}()

	return nil
}

func RenameGoDirectory(goVersionName string) error {
	target := fmt.Sprintf("%s/%s", getVersionsDir(), "go")

	if err := os.Rename(target, fmt.Sprintf("%s/%s", getVersionsDir(), goVersionName)); err != nil {
		return err
	}

	return nil
}

func RemoveTarFile() error {
	err := os.Remove(getTarFile())
	if err != nil {
		return err
	}
	return nil
}

func CreateExecutableSymlink(goVersionName string) error {
	target := fmt.Sprintf("%s/%s/bin", getVersionsDir(), goVersionName)

	files, err := os.ReadDir(target)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		newFile := fmt.Sprintf("%s/%s", target, f.Name())
		link := fmt.Sprintf("%s/%s", getBinDir(), f.Name())

		if _, err := os.Lstat(link); err == nil {
			os.Remove(link)
		}

		if err := os.Symlink(newFile, link); err != nil {
			return err
		}
		os.Chmod(link, 0700)
	}

	return nil
}

func UpdateRecentVersion(goVersionName string) error {
	path := getCurrentVersionFile()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, goVersionName)
	return err
}

func GetRecentVersion() string {
	path := getCurrentVersionFile()

	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func VersionExists(goVersion string) bool {
	target := getVersionsDir()

	if _, err := os.Stat(fmt.Sprintf("%s/%s", target, goVersion)); !os.IsNotExist(err) {
		return true
	}

	return false
}
