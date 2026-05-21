package implementations

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Logger struct {
	file *os.File

	mu   sync.Mutex
	logs []string

	maxLogs int
}

func New(path string) (*Logger, error) {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		file:    file,
		maxLogs: 500,
	}

	mw := io.MultiWriter(file, l)

	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return l, nil
}

func (l *Logger) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, msg)

	if len(l.logs) > l.maxLogs {
		l.logs = l.logs[len(l.logs)-l.maxLogs:]
	}

	return len(p), nil
}

func (l *Logger) Logs() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]string, len(l.logs))
	copy(out, l.logs)

	return out
}

func (l *Logger) Close() error {
	return l.file.Close()
}
