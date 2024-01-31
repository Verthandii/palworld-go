package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger(prefix string) *Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.SetPrefix(fmt.Sprintf("【%s】", prefix))
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}
