package logger

import (
	"fmt"
	"log"
	"sync"
)

// Logger is a wrapper around the standard [log.Logger] that adds a prefix to each
// loger instance.
type Logger struct {
	*log.Logger
}

const (
	prefixBrackets       = "[%s]"
	prefixFormat         = "% -*s | "
	prefixBracketsLength = len(prefixBrackets) - 2
)

var (
	loggers      = make(map[string]*Logger)
	loggersMux   sync.Mutex
	prefixLength int
)

// New creates a new logger with the given prefix that writes to the standard
// logger destination.
func New(prefix string) *Logger {
	loggersMux.Lock()
	defer loggersMux.Unlock()

	if l, ok := loggers[prefix]; ok {
		return l
	}

	if len(prefix) > prefixLength {
		prefixLength = len(prefix)
		updatePrefixes()
	}

	loggers[prefix] = &Logger{log.New(
		log.Writer(),
		getPrefix(prefix),
		log.LstdFlags|log.Lmsgprefix,
	)}
	return loggers[prefix]
}

func updatePrefixes() {
	for prefix, l := range loggers {
		l.SetPrefix(getPrefix(prefix))
	}
}

func getPrefix(prefix string) string {
	prefix = fmt.Sprintf(prefixBrackets, prefix)
	return fmt.Sprintf(prefixFormat, prefixBracketsLength+prefixLength, prefix)
}
