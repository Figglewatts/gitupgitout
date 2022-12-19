package logging

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Verbose(format string, v ...any)
	Print(format string, v ...any)
	Error(format string, v ...any)
}

type StdLogger struct {
	LogVerbose bool

	log *log.Logger
	err *log.Logger
}

func NewStd() StdLogger {
	return StdLogger{
		false,
		log.New(os.Stdout, "", log.LstdFlags),
		log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

func (l StdLogger) Verbose(format string, v ...any) {
	if !l.LogVerbose {
		return
	}
	l.log.Print(l.msg(format, v...))
}

func (l StdLogger) Print(format string, v ...any) {
	l.log.Print(l.msg(format, v...))
}

func (l StdLogger) Error(format string, v ...any) {
	l.err.Print(l.msg(format, v...))
}

func (l StdLogger) msg(format string, v ...any) string {
	return fmt.Sprintf(format+"\n", v...)
}
