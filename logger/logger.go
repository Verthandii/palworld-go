package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		prefix: fmt.Sprintf("【%s】", prefix),
	}
}

func (l *Logger) Println(v ...interface{}) {
	v = append([]interface{}{l.prefix}, v...)
	l.logger.Println(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(l.prefix+format, v...)
}
