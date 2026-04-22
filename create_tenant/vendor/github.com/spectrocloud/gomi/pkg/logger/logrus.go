package logger

import (
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Entry
}

var logrusApplog *logrus.Logger

func newLogrusLogger(config ApplogConfig) Logger {
	if config.LogLevel == "" {
		config.LogLevel = DEFAULT
	}
	level, _ := logrus.ParseLevel(config.LogLevel)

	var formatter logrus.Formatter
	if config.JSONFormatter {
		formatter = &logrus.JSONFormatter{}
	} else {
		formatter = &logrus.TextFormatter{
			DisableLevelTruncation: true,
		}
	}

	loggerConfig := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Level:     level,
	}

	if config.ReportCaller {
		reportCaller = true
	}

	logrusApplog = loggerConfig
	logrusEntry := logrusLogger{
		logger: logrus.NewEntry(logrusApplog),
	}
	return &logrusEntry
}

func (l *logrusLogger) WithFields(f Fields) Logger {
	fields := logrus.Fields{}
	for i, val := range f {
		fields[i] = val
	}

	fieldLogger := logrus.NewEntry(logrusApplog)
	fieldLogger = fieldLogger.WithFields(fields)

	return &logrusLogger{
		logger: fieldLogger,
	}
}

// SetLevel updates the log level
func (l *logrusLogger) SetLogLevel(lvl string) {
	level, _ := logrus.ParseLevel(lvl)
	l.logger.Logger.SetLevel(level)
}

func (l *logrusLogger) SetFormatter(formatter string) {
	switch formatter {
	case "text":
		f := &logrus.TextFormatter{
			DisableLevelTruncation: true,
		}
		l.logger.Logger.SetFormatter(f)
	case "json":
		f := &logrus.JSONFormatter{}
		l.logger.Logger.SetFormatter(f)
	}
}

func (l *logrusLogger) getReportCallerEntry() *logrus.Entry {
	var fun string
	var skip = minCallerSkip
	if !reflect.DeepEqual(applog, l) {
		skip = minCallerSkip - 1
	}

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "unknown"
		line = 0
		fun = "unknown"
	}

	if fun != "unknown" {
		fun = runtime.FuncForPC(pc).Name()
	}

	// Create a new entry with the caller fields
	entry := l.logger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
		"func": fun,
	})

	return entry
}

func (l *logrusLogger) Panic(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Panicf(format, args...)
}

func (l *logrusLogger) Fatal(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Fatalf(format, args...)
}

func (l *logrusLogger) Error(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Errorf(format, args...)
}

func (l *logrusLogger) Warn(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Warnf(format, args...)
}

func (l *logrusLogger) Info(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Infof(format, args...)
}

func (l *logrusLogger) Debug(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Debugf(format, args...)
}

func (l *logrusLogger) Trace(format string, args ...interface{}) {
	entry := l.logger
	if reportCaller {
		entry = l.getReportCallerEntry()
	}
	entry.Tracef(format, args...)
}
