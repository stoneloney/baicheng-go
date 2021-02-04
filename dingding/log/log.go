package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	LogLevelNull    = 0
	LogLevelTrace   = 1
	LogLevelDebug   = 2
	LogLevelInfo    = 3
	LogLevelWarning = 4
	LogLevelError   = 5
	LogLevelFatal   = 6
)

var DefaultLogLevel int = LogLevelDebug

type Conf struct {
	Name    string
	MaxSize int64
	MaxNum  int
}

type Logger struct {
	Writer io.Writer
	Level  int
}

func NewLogger() *Logger {
	conf := Conf{
		Name:    "app",      // 日志名称
		MaxSize: 1073741824, // 1G
		MaxNum:  2,          // 保留2个日志文件
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	realPath := fmt.Sprintf("%s/%s", dir, "log")
	_, err = os.Stat(realPath)
	if os.IsNotExist(err) {
		_ = os.Mkdir(realPath, os.ModePerm)
	}
	logWriter, err := NewFileWriter(realPath, conf.Name, conf.MaxSize, conf.MaxNum)
	if err != nil {
		panic(err)
	}
	return &Logger{
		Writer: logWriter,
		Level:  DefaultLogLevel,
	}
}

func (l *Logger) SetLevel(level int) {
	l.Level = level
}

func (l *Logger) GetLevel() int {
	return l.Level
}

func (l *Logger) Debug(mess string) {
	if l.Level > LogLevelDebug {
		return
	}
	log.Output(2, string("[DEBUG] ")+mess)
}

func (l *Logger) Info(mess string) {
	if l.Level > LogLevelInfo {
		return
	}
	log.Output(2, string("[INFO] ")+mess)
}

func (l *Logger) Warning(mess string) {
	if l.Level > LogLevelWarning {
		return
	}
	log.Output(2, string("[WARNING] ")+mess)
}

func (l *Logger) Error(mess string) {
	if l.Level > LogLevelError {
		return
	}
	log.Output(2, string("[ERROR] ")+mess)
}

func (l *Logger) Fatal(mess string) {
	if l.Level > LogLevelFatal {
		return
	}
	log.Output(2, string("[Fatal] ")+mess)
}

func (l *Logger) Panic(mess string) {
	log.Output(2, string("[Panic]")+mess)
}
