package testutils

type FakeStdout struct {
	messages []string
	Closed   bool
}

func (s FakeStdout) GetPrintMessages() []string {
	return s.messages
}

func (s *FakeStdout) Write(p []byte) (n int, err error) {
	s.messages = append(s.messages, string(p))

	return 0, nil
}

func (s *FakeStdout) Close() error {
	s.Closed = true
	s.messages = []string{}
	return nil
}
