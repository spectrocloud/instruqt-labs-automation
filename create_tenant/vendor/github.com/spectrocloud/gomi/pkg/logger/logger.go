package logger

import "context"

// Global application logger
var applog Logger

type key int

// LoggerKey is the key to be used for logger context
const LoggerKey key = 1

const (
	PANIC   = "panic"
	FATAL   = "fatal"
	ERROR   = "error"
	WARN    = "warn"
	INFO    = "info"
	DEBUG   = "debug"
	TRACE   = "trace"
	DEFAULT = "info"
)

// Fields map to be used with WithFields
type Fields map[string]interface{}

var (
	minCallerSkip = 3
	reportCaller  = true
)

// ApplogConfig sets the config for the logger
type ApplogConfig struct {
	// Logs file name, function name and line number if set
	// Setting this adds measurable overhead
	ReportCaller bool

	// Logs in JSON format if set.
	// Default is text format
	JSONFormatter bool

	// Log level for messages.
	// Default level is info
	LogLevel string
}

// Logger holds the exposed logger functionality
type Logger interface {
	Panic(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	SetLogLevel(lvl string)
	WithFields(fields Fields) Logger
	SetFormatter(formatter string)
}

func init() {
	defaultConfig := ApplogConfig{}
	InitApplog(defaultConfig)
}

// InitApplog initializes the application logger
func InitApplog(config ApplogConfig) {
	logger := newLogrusLogger(config)
	applog = logger
}

func WithFields(fields Fields) Logger {
	return applog.WithFields(fields)
}

func WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return applog
	}

	if ctxLogger, ok := ctx.Value(LoggerKey).(Logger); ok {
		return ctxLogger
	}
	return applog
}

func SetContext(ctx context.Context, fields Fields) context.Context {
	ctxLogger := WithFields(fields)
	return context.WithValue(ctx, LoggerKey, ctxLogger)
}

func SetLogLevel(lvl string) {
	applog.SetLogLevel(lvl)
}

func SetTextFormatter() {
	applog.SetFormatter("text")
}

func SetJSONFormatter() {
	applog.SetFormatter("json")
}

func SetReportCaller() {
	reportCaller = true
}

func UnsetReportCaller() {
	reportCaller = false
}

func Panic(format string, args ...interface{}) {
	applog.Panic(format, args...)
}

func Fatal(format string, args ...interface{}) {
	applog.Fatal(format, args...)
}

func Error(format string, args ...interface{}) {
	applog.Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	applog.Warn(format, args...)
}

func Info(format string, args ...interface{}) {
	applog.Info(format, args...)
}

func Debug(format string, args ...interface{}) {
	applog.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	applog.Trace(format, args...)
}
