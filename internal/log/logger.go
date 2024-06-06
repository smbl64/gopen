package log

import (
	stdlog "log"
	"os"
)

type Logger struct {
	l     *stdlog.Logger
	debug bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{
		debug: debug,
		l:     stdlog.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Infof(format string, v ...any) {
	l.l.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	if !l.debug {
		return
	}

	l.l.Printf(format, v...)
}
