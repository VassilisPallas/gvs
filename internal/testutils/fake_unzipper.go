package testutils

type FakeUnzipper struct{}

func (FakeUnzipper) ExtractTarSource(dst string, src string) error { return nil }

func (FakeUnzipper) ExtractZipSource(dst string, src string) error { return nil }
