// Package logger provides an interface to print logs
// The implementation is both for cli messages to the user
// as well as for an actual logger that can be used for
// debugging purposes
package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// Logger is the interface that wraps the basic methods for printing messages to the user
// and logging information.
//
// PrintMessage prints the given message to the cli. It is similar to Printf function from the fmt package,
// which means it accepts a format specifier and the variables to be printed.
// This method could also have the freedom to print the message to the logger.
//
// PrintError is similar function to PrintMessage, but it is used to differentiate when information and error
// messages as different methods. This method could also have the freedom to print the error message to the logger.
//
// Info is the method that must be used for Info logger messages. It is similar to Printf function from the fmt package,
// which means it accepts a format specifier and the variables to be printed.
//
// Error is the method that must be used for Error logger messages. It is similar to Printf function from the fmt package,
// which means it accepts a format specifier and the variables to be printed.
//
// SetLogWriter specified the output destination for the logger.
//
// Close closed the logWriter instance that is passed to the method SetLogWriter (if any).
type Logger interface {
	PrintMessage(format string, a ...any)
	PrintError(format string, a ...any)
	Info(format string, a ...any)
	Error(format string, a ...any)
	SetLogWriter(logWriter io.WriteCloser)
	Close()
}

// Log is the struct that implements the Logger interface
//
// Go struct accepts three fields, the cliWriter which is the output destination for the cli messages,
// the logWriter which is the output destination for the log messages, and the logger, which is a
// *log.Logger instance.
type Log struct {
	cliWriter io.Writer
	logWriter io.WriteCloser
	logger    *log.Logger
}

// addNewLine adds a new line to the precified format specifier
// if it does not contain one.
// This method is used to always make sure that both cli and log messages are nicely printed.
func (l Log) addNewLine(format string) string {
	newLine := "\n"

	if strings.HasSuffix(format, newLine) {
		return format
	}

	return format + newLine
}

// PrintMessage sends the given message both to the cli and to the logger outputs
//
// If cliWriter does not contain a non-null value, then the method returns
// without outputting anything.
//
// To output the message to the logger, is it using the Info method from the Logger interface.
func (l *Log) PrintMessage(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	message := l.addNewLine(format)

	fmt.Fprintf(l.cliWriter, message, a...)
	l.Info(message, a...)
}

// PrintMessage sends the given error message both to the cli and to the logger outputs
//
// If cliWriter does not contain a non-null value, then the method returns
// without outputting anything.
//
// To output the error message to the logger, is it using the Error method from the Logger interface.
func (l *Log) PrintError(format string, a ...any) {
	if l.cliWriter == nil {
		return
	}

	message := l.addNewLine(format)

	fmt.Fprintf(l.cliWriter, l.addNewLine(message), a...)
	l.Error(message, a...)
}

// Info sends the given message to the logger output
//
// If the logWriter does not contain a non-null value, then the method returns
// without outputting anything.
//
// It is also adding the prefix `INFO: ` to the logger instance
func (l *Log) Info(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	message := l.addNewLine(format)

	l.logger.SetPrefix("INFO: ")
	l.logger.Printf(message, a...)
}

// Error sends the given error message to the logger output
//
// If the logWriter does not contain a non-null value, then the method returns
// without outputting anything.
//
// It is also adding the prefix `ERROR: ` to the logger instance
func (l *Log) Error(format string, a ...any) {
	if l.logWriter == nil {
		return
	}

	message := l.addNewLine(format)

	l.logger.SetPrefix("ERROR: ")
	l.logger.Printf(message, a...)
}

// Close closed the logger output writer
//
// If the logWriter does not contain a non-null value, then the method returns
// without closing anything.
func (l *Log) Close() {
	if l.logWriter == nil {
		return
	}

	l.logWriter.Close()
}

// SetLogWriter sets the logger output writer
//
// It is storing it to both the struct so it can be used later,
// and also in the logger as an output writer
//
// If the logWriter is non-null value to the struct,
// it returns without re-storing it.
func (l *Log) SetLogWriter(logWriter io.WriteCloser) {
	if l.logWriter != nil {
		return
	}

	l.logWriter = logWriter
	l.logger.SetOutput(logWriter)
}

// New returns a *Log instance that implements the Logger interface.
// Each call to New returns a distinct *Log instance even if the parameters are identical.
func New(cliWriter io.Writer, logWriter io.WriteCloser) *Log {
	customLogger := log.New(logWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	return &Log{cliWriter: cliWriter, logWriter: logWriter, logger: customLogger}
}
