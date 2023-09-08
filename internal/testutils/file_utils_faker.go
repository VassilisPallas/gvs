package testutils

type FakeFiler struct {
	TarFile            string
	VersionDir         string
	CurrentVersionFile string
}

func (FakeFiler) CreateInitFiles() error {
	return nil
}

func (FakeFiler) GetAppDir() string {
	return ""
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

func (FakeFiler) GetVersionResponseFile() string {
	return ""
}

func (FakeFiler) GetHomeDirectory() string {
	return ""
}
