# Logging

## Using this package

### Importing logger

```go
log "github.com/spectrocloud/gomi/pkg/logger"
```

### Default configuration
* Formatter: text
* LogLevel: Info
* ReportCaller: true

### Initializing logger

```go
config := log.ApplogConfig{
    JSONFormatter: true,
    ReportCaller:  false,
    LogLevel:      log.TRACE,
}
log.InitApplog(config)
```

### Log levels

```go
log.Panic("Panic level") // Logs and panics
log.Fatal("Fatal level") // Logs and exits
log.Error("Error level")
log.Warn("Warn level")
log.Info("Info level")
log.Debug("Debug level")
log.Trace("Trace level")
```

### Set LogLevel

```go
log.SetLogLevel(log.DEBUG)
log.Debug("Debug level set")
log.Trace("Wont log trace")
```

### WithField logger

```go
fieldLogger := log.WithFields(log.Fields{
    "id":   "1234",
    "name": "Alice",
})
fieldLogger.Info("Logged with fields")
```

### Context logger

```go
// Set context
ctx := log.SetContext(context.Background(), log.Fields{"reqID": "1001", "tenantID": "2001"})

// Get context
ctxLogger := log.WithContext(ctx)
ctxLogger.Info("Logged with context")
```

### Change formatter

```go
// For text formatter
log.SetTextFormatter()
log.Info("Logged in text format")

// For JSON formatter
log.SetJSONFormatter()
log.Info("Logged in json format")
```

### Set ReportCaller

```go
// Set reporting caller function
log.SetReportCaller()
log.Info("Will log function name, file name and line number")

// Unset reporting caller function
log.UnsetReportCaller()
```
