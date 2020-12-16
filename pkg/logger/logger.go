package logger

import (
	"fmt"
	"log"
)

type Interface interface {
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
}

type Logger struct {
	IsDebug bool
}

// Logger is awful, not for production, made just for dev necessity
func NewLogger(isDebug bool) *Logger {
	return &Logger{IsDebug: isDebug}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if !l.IsDebug {
		return
	}

	log.Println("[DEBUG] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	log.Println("[INFO] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	log.Println("[WARN] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Println("[ERROR] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	log.Fatalln("[FATAL] " + fmt.Sprintf(format, v...))
}
