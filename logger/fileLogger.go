package logger

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type FileLogger struct {
	level         LogLevel
	logFilePath   string
	logFileName   string
	maxFileSize   int64
	fileWriter    *os.File
	errFileWriter *os.File
	logChan       chan *logPack
	consoleLogger *Logger
}
type logPack struct {
	level    LogLevel
	time     string
	fileName string
	funcName string
	line     int
	msg      string
	file     *os.File
	isError  bool
}

func NewFileLogger(level LogLevel, path, name string, maxSize int64, isConsole bool) *FileLogger {
	logger := &FileLogger{
		level:       level,
		logFilePath: path,
		logFileName: name,
		maxFileSize: maxSize,
		logChan:     make(chan *logPack, 50000),
	}
	if isConsole {
		logger.consoleLogger = NewLogger(level)
	} else {
		logger.consoleLogger = nil
	}
	logger.initWriter()
	go logger.writeLog()
	return logger
}

func (f *FileLogger) Debug(format string, a ...any) {
	f.createInfo(L_DEBUG, format, a...)
}
func (f *FileLogger) Warning(format string, a ...any) {
	f.createInfo(L_WARNING, format, a...)
}
func (f *FileLogger) Error(format string, a ...any) {
	f.createInfo(L_ERROR, format, a...)
}
func (f *FileLogger) Error1(err error) {
	f.createInfo(L_ERROR, err.Error())
}

func (f *FileLogger) createInfo(targetLogLevel LogLevel, format string, a ...any) {
	if f.level <= targetLogLevel {
		f.makeAndWriteInfo(false, f.fileWriter, targetLogLevel, format, a...)
	}
	//当错误等级大于等于ERROR时，需要将错误单独写到一个文件中去
	if targetLogLevel >= L_ERROR {
		f.makeAndWriteInfo(true, f.errFileWriter, targetLogLevel, format, a...)
	}
}
func (f *FileLogger) makeAndWriteInfo(isError bool, file *os.File, targetLogLevel LogLevel, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fileName, funcName, line := getLogInfo(5)
	pack := &logPack{
		level:    targetLogLevel,
		time:     time.Now().Format("2006-01-02 15:04:05"),
		fileName: fileName,
		funcName: funcName,
		line:     line,
		msg:      msg,
		file:     file,
		isError:  isError,
	}
	select {
	case f.logChan <- pack:
	default:
	}
}
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		panic("checkSize Error:")
	}
	size := fileInfo.Size()
	return size >= f.maxFileSize
}
func (f *FileLogger) writeLog() {
	for true {
		select {
		case pack := <-f.logChan:
			file := pack.file
			if f.checkSize(file) {
				//需要切割

				//1.关闭旧的文件写入
				oldPath := path.Join(f.logFilePath, file.Name())
				newPath := oldPath + ".ok" + strconv.FormatInt(time.Now().UnixNano(), 10)
				err := file.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
				//2.重命名旧的文件
				err = os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Println(err)
					return
				}
				//3.创建新的文件写入
				newFile, err := os.OpenFile(oldPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
				//4.设置新的file
				if *file == *f.fileWriter {
					f.fileWriter = newFile
				} else {
					f.errFileWriter = newFile
				}
			}
			msgFormat := fmt.Sprintf("[%s][%s][%s:%s:%d]:%s\n",
				getLogLevel(pack.level),
				pack.time,
				pack.fileName,
				pack.funcName,
				pack.line,
				pack.msg)
			if f.consoleLogger != nil && !pack.isError {
				f.consoleLogger.CreateInfo(pack.level, pack.fileName, pack.line, pack.msg)
			}
			_, err := fmt.Fprintf(file, msgFormat)
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
func (f *FileLogger) initWriter() {
	p := path.Join(f.logFilePath, f.logFileName)
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	errFile, err := os.OpenFile(p+".err", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.fileWriter = file
	f.errFileWriter = errFile
}
func (f *FileLogger) Close() {
	err := f.fileWriter.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = f.errFileWriter.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.fileWriter = nil
	f.errFileWriter = nil
}
