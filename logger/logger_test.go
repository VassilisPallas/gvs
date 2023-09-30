package logger_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/VassilisPallas/gvs/internal/testutils"
	"github.com/VassilisPallas/gvs/logger"
	"github.com/google/go-cmp/cmp"
)

func TestPrintMessage(t *testing.T) {
	msg := "some message"
	cliWriter := &testutils.FakeStdout{}
	logWriter := &testutils.FakeStdout{}

	log := logger.New(cliWriter, logWriter)

	log.PrintMessage(msg)

	printedMessages := cliWriter.GetPrintMessages()
	expectedPrintedMessages := []string{msg}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong print messages received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}

	logMessage := logWriter.GetPrintMessages()[0]
	regex := fmt.Sprintf("(INFO:) (\\d{4}\\/\\d{2}\\/\\d{2}) (\\d{2}:\\d{2}:([0-9]*[.])?[0-9]+) .+: %s", msg)
	match, _ := regexp.MatchString(regex, logMessage)
	if !match {
		t.Errorf("Wrong log messages received, got=%s", logMessage)
	}
}

func TestPrintError(t *testing.T) {
	msg := "some error"
	cliWriter := &testutils.FakeStdout{}
	logWriter := &testutils.FakeStdout{}

	log := logger.New(cliWriter, logWriter)

	log.PrintError(msg)

	printedMessages := cliWriter.GetPrintMessages()
	expectedPrintedMessages := []string{msg}
	if !cmp.Equal(printedMessages, expectedPrintedMessages) {
		t.Errorf("Wrong print messages received, got=%s", cmp.Diff(expectedPrintedMessages, printedMessages))
	}

	logMessage := logWriter.GetPrintMessages()[0]
	regex := fmt.Sprintf("(ERROR:) (\\d{4}\\/\\d{2}\\/\\d{2}) (\\d{2}:\\d{2}:([0-9]*[.])?[0-9]+) .+: %s", msg)
	match, _ := regexp.MatchString(regex, logMessage)
	if !match {
		t.Errorf("Wrong log messages received, got=%s", logMessage)
	}
}

func TestClose(t *testing.T) {
	cliWriter := &testutils.FakeStdout{}
	logWriter := &testutils.FakeStdout{}

	log := logger.New(cliWriter, logWriter)

	if logWriter.Closed != false {
		t.Error("logWriter should not be closed")
	}

	log.Close()

	if logWriter.Closed != true {
		t.Error("logWriter should be closed")
	}
}
