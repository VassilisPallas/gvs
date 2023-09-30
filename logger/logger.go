package logger

import (
	"fmt"
	"io"
	"log"
)

type Logger interface {
	PrintMessage(format string, a ...any)
	PrintError(format string, a ...any)
	Info(format string, a ...any)
	Error(format string, a ...any)
	Close()
}

type Log struct {
	cliWriter io.Writer
	logWriter io.WriteCloser
	logger    *log.Logger
}

func (l *Log) PrintMessage(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	fmt.Fprintf(l.cliWriter, format, a...)
	l.Info(format, a...)
}

func (l *Log) PrintError(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	fmt.Fprintf(l.cliWriter, format, a...)
	l.Error(format, a...)
}

func (l *Log) Info(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	l.logger.SetPrefix("INFO: ")
	l.logger.Printf(format, a...)
}

func (l *Log) Error(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	l.logger.SetPrefix("ERROR: ")
	l.logger.Printf(format, a...)
}

func (l *Log) Close() {
	l.logWriter.Close()
}

func New(cliWriter io.Writer, logWriter io.WriteCloser) *Log {
	customLogger := log.New(logWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	return &Log{cliWriter: cliWriter, logWriter: logWriter, logger: customLogger}
}
