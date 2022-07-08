package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level: level,
	}
}

func (l *Logger) CreateInfo(targetLogLevel LogLevel, fileName string, line int, msg string) {
	if l.level <= targetLogLevel {
		fmt.Printf("[%s][%s][%s:%d]:%s\n",
			getLogLevel(targetLogLevel),
			time.Now().Format("2006-01-02 15:04:05"),
			fileName,
			line,
			msg)
	}
}
