package testutils

type FakeFiler struct {
	TarFile                 string
	VersionDir              string
	CurrentVersionFile      string
	AppDir                  string
	VersionResponseFileName string
}

func (FakeFiler) CreateInitFiles() error {
	return nil
}

func (ff FakeFiler) GetAppDir() string {
	return ff.AppDir
}

func (ff FakeFiler) GetVersionsDir() string {
	return ff.VersionDir
}

func (ff FakeFiler) GetTarFile() string {
	return ff.TarFile
}

func (FakeFiler) GetBinDir() string {
	return ""
}

func (ff FakeFiler) GetCurrentVersionFile() string {
	return ff.CurrentVersionFile
}

func (ff FakeFiler) GetVersionResponseFile() string {
	return ff.VersionResponseFileName
}

func (FakeFiler) GetHomeDirectory() string {
	return ""
}
