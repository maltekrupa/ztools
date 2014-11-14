package zlog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	prefix string
}

type LogLevel uint8

const (
	prefixFormat = "%s [%s] %s: "
)

const (
	LOG_FATAL LogLevel = iota
	LOG_ERROR LogLevel = iota
	LOG_WARN  LogLevel = iota
	LOG_INFO  LogLevel = iota
	LOG_DEBUG LogLevel = iota
	LOG_TRACE LogLevel = iota
)

var (
	levelNames = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
)

func (level LogLevel) String() string {
	if level > LOG_TRACE {
		level = LOG_TRACE
	}
	return levelNames[level]
}

var defaultLogger = New(os.Stdout, "log")

func New(out io.Writer, prefix string) *Logger {
	logger := Logger{
		out:    out,
		prefix: prefix,
	}
	return &logger
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.doPrint(LOG_FATAL, v...)
	os.Exit(1)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.doPrintf(LOG_FATAL, format, v...)
	os.Exit(1)
}

func (logger *Logger) Error(v ...interface{}) {
	logger.doPrint(LOG_ERROR, v...)
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.doPrintf(LOG_ERROR, format, v...)
}

func (logger *Logger) Warn(v ...interface{}) {
	logger.doPrint(LOG_WARN, v...)
}

func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.doPrintf(LOG_WARN, format, v...)
}

func (logger *Logger) Info(v ...interface{}) {
	logger.doPrint(LOG_INFO, v...)
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.doPrintf(LOG_INFO, format, v...)
}

func (logger *Logger) Debug(v ...interface{}) {
	logger.doPrint(LOG_DEBUG, v...)
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.doPrintf(LOG_DEBUG, format, v...)
}

func (logger *Logger) Trace(v ...interface{}) {
	logger.doPrint(LOG_TRACE, v...)
}

func (logger *Logger) Tracef(format string, v ...interface{}) {
	logger.doPrintf(LOG_TRACE, format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Trace(v ...interface{}) {
	defaultLogger.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

func (logger *Logger) Print(level LogLevel, v ...interface{}) {
	logger.doPrint(level, v...)
	if level == LOG_FATAL {
		os.Exit(1)
	}
}

func (logger *Logger) Printf(level LogLevel, format string, v ...interface{}) {
	logger.doPrintf(level, format, v...)
	if level == LOG_FATAL {
		os.Exit(1)
	}
}

func Print(level LogLevel, v ...interface{}) {
	defaultLogger.Print(level, v...)
}

func Printf(level LogLevel, format string, v ...interface{}) {
	defaultLogger.Printf(level, format, v...)
}

func (logger *Logger) doPrint(level LogLevel, v ...interface{}) {
	timestamp := time.Now().Format(time.StampMilli)
	s := fmt.Sprintf(prefixFormat, timestamp, level.String(), logger.prefix)
	front := []byte(s)
	back := []byte(fmt.Sprint(v...))
	logger.writeLine(front, back)
}

func (logger *Logger) doPrintf(level LogLevel, format string, v ...interface{}) {
	timestamp := time.Now().Format(time.StampMilli)
	s := fmt.Sprintf(prefixFormat, timestamp, level.String(), logger.prefix)
	front := []byte(s)
	back := []byte(fmt.Sprintf(format, v...))
	logger.writeLine(front, back)
}

func (logger *Logger) writeLine(b ...[]byte) {
	logger.mu.Lock()
	for _, chunk := range b {
		logger.out.Write(chunk)
	}
	logger.out.Write([]byte{'\n'})
	logger.mu.Unlock()
}
