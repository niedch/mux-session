package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	enabled bool
	output  io.Writer
}

var defaultLogger = New(false, os.Stderr)

func New(enabled bool, output io.Writer) *Logger {
	if output == nil {
		output = os.Stderr
	}
	return &Logger{
		enabled: enabled,
		output:  output,
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.enabled {
		fmt.Fprintf(l.output, format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.enabled {
		fmt.Fprintf(l.output, format, v...)
	}
	os.Exit(1)
}

func (l *Logger) Enabled() bool {
	return l.enabled
}

func (l *Logger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func SetDefault(l *Logger) {
	defaultLogger = l
}

func GetDefault() *Logger {
	return defaultLogger
}

func Printf(format string, v ...interface{}) {
	defaultLogger.Printf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func SetEnabled(enabled bool) {
	defaultLogger.SetEnabled(enabled)
}

func Enabled() bool {
	return defaultLogger.Enabled()
}
