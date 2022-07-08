package logger

import (
	"fmt"
	"path"
	"runtime"
)

type LogLevel uint8

const (
	L_DEBUG LogLevel = iota
	L_WARNING
	L_ERROR
)

func getLogLevel(level LogLevel) string {
	switch level {
	case L_DEBUG:
		return "DEBUG"
	case L_WARNING:
		return "WARNING"
	case L_ERROR:
		return "ERROR"
	default:
		return "nil"
	}
}

func getLogInfo(skip int) (fileName, funcName string, l int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller Error")
		return
	}
	fileName = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	l = line
	return
}
