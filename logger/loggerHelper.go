package logger

var loggerMgr *FileLogger

func init() {
	loggerMgr = NewFileLogger(L_DEBUG,
		"ServerLogger", "GameServerInfoLogger.txt", 1024*1024, true)
}

func PopError(err error) {
	loggerMgr.Error1(err)
}
func PopErrorInfo(info string, a ...any) {
	loggerMgr.Error(info, a...)
}
func PopDebug(msg string, a ...any) {
	loggerMgr.Debug(msg, a...)
}
func PopWarning(msg string, a ...any) {
	loggerMgr.Warning(msg, a...)
}
