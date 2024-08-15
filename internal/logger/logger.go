package logger

import (
	"log"
	"os"
)

// A simple logger with outputs to stdout and stderr.
type Logger struct {
	debug   *log.Logger
	err     *log.Logger
	verbose bool
}

func New(verbose bool) *Logger {
	return &Logger{
		debug:   log.New(os.Stdout, "", 0),
		err:     log.New(os.Stderr, "", 0),
		verbose: verbose,
	}
}

func (l *Logger) Debug(msg string) {
	if l.verbose {
		l.debug.Println(msg)
	}
}

func (l *Logger) Error(msg string) {
	l.err.Println(msg)
}

func (l *Logger) Fatal(msg string) {
	l.err.Println(msg)
	os.Exit(1)
}
