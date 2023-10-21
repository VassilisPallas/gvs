package testutils

type FakeUnzipper struct {
	ExtractTarSourceError error
}

func (u FakeUnzipper) ExtractTarSource(dst string, src string) error { return u.ExtractTarSourceError }

func (FakeUnzipper) ExtractZipSource(dst string, src string) error { return nil }
