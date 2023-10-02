package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Logger interface {
	PrintMessage(format string, a ...any)
	PrintError(format string, a ...any)
	Info(format string, a ...any)
	Error(format string, a ...any)
	Close()
	SetLogWriter(logWriter io.WriteCloser)
}

type Log struct {
	cliWriter io.Writer
	logWriter io.WriteCloser
	logger    *log.Logger
}

func (l Log) addNewLine(format string) string {
	newLine := "\n"

	if strings.HasSuffix(format, newLine) {
		return format
	}

	return format + newLine
}

func (l *Log) PrintMessage(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	message := l.addNewLine(format)

	fmt.Fprintf(l.cliWriter, message, a...)
	l.Info(message, a...)
}

func (l *Log) PrintError(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	message := l.addNewLine(format)

	fmt.Fprintf(l.cliWriter, l.addNewLine(message), a...)
	l.Error(message, a...)
}

func (l *Log) Info(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	message := l.addNewLine(format)

	l.logger.SetPrefix("INFO: ")
	l.logger.Printf(message, a...)
}

func (l *Log) Error(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	message := l.addNewLine(format)

	l.logger.SetPrefix("ERROR: ")
	l.logger.Printf(message, a...)
}

func (l *Log) Close() {
	l.logWriter.Close()
}

func (l *Log) SetLogWriter(logWriter io.WriteCloser) {
	l.logWriter = logWriter
	l.logger.SetOutput(logWriter)
}

func New(cliWriter io.Writer, logWriter io.WriteCloser) *Log {
	customLogger := log.New(logWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	return &Log{cliWriter: cliWriter, logWriter: logWriter, logger: customLogger}
}
