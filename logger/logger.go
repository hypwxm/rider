package logger

import (
	"log"
	"fmt"
)

const (
	LevelError = iota
	LevelWarning
	LevelInformational
	LevelDebug
)

type Logger struct {
	level int
	error *log.Logger
	warn *log.Logger
	info *log.Logger
	debug *log.Logger
}


func (l *Logger) Error(format string, v ...interface{}) {
	if LevelError > l.level {
		return
	}

	msg := fmt.Sprintf("[ERROR] " + format, v...)
	l.error.Printf(msg)
}

//设置日志级别
func (l *Logger) SetLevel(i int) {
	l.level = i
}