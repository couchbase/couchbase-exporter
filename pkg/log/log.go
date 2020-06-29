package log

import (
	gofmt "fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	timestampFormat = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05.000Z07:00",
	)
	defaultLogger = New(&Config{})
)

type logger struct {
	debug log.Logger
	info  log.Logger
	warn  log.Logger
	err   log.Logger

	lvl    string
	format string
}

func (l *logger) Debug(fmt string, v ...interface{}) error {
	return l.debug.Log(prepareLog(fmt, v...))
}

func (l *logger) Info(fmt string, v ...interface{}) error {
	return l.info.Log(prepareLog(fmt, v...))
}

func (l *logger) Warn(fmt string, v ...interface{}) error {
	return l.warn.Log(prepareLog(fmt, v...))
}

func (l *logger) Error(fmt string, v ...interface{}) error {
	return l.err.Log(prepareLog(fmt, v...))
}

func (l *logger) SetLevel(lvl string) {
	*l = *(New(&Config{Level: lvl, Format: l.format}))
}

func (l *logger) SetFormat(format string) {
	*l = *(New(&Config{Level: l.lvl, Format: format}))
}

func Debug(fmt string, v ...interface{}) error {
	return defaultLogger.Debug(fmt, v...)
}

func Info(fmt string, v ...interface{}) error {
	return defaultLogger.Info(fmt, v...)
}

func Warn(fmt string, v ...interface{}) error {
	return defaultLogger.Warn(fmt, v...)
}

func Error(fmt string, v ...interface{}) error {
	return defaultLogger.Error(fmt, v...)
}

func SetLevel(lvl string) {
	*defaultLogger = *(New(&Config{Level: lvl, Format: defaultLogger.format}))
}

func SetFormat(format string) {
	*defaultLogger = *(New(&Config{Level: defaultLogger.lvl, Format: format}))
}

type Config struct {
	Level  string
	Format string
}

func New(cfg *Config) *logger {
	var l log.Logger
	if cfg.Format == "json" {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}
	l = log.With(level.NewFilter(l, parseLevel(cfg.Level)), "date", timestampFormat, "caller", log.DefaultCaller)
	return &logger{debug: level.Debug(l), info: level.Info(l), warn: level.Warn(l), err: level.Error(l), lvl: cfg.Level, format: cfg.Format}
}

func prepareLog(fmt string, v ...interface{}) (string, string) {
	return "message", gofmt.Sprintf(fmt, v...)
}

func parseLevel(l string) level.Option {
	lvl := level.AllowWarn()
	switch l {
	case "debug":
		lvl = level.AllowDebug()
	case "info":
		lvl = level.AllowInfo()
	case "warn", "warning":
		lvl = level.AllowWarn()
	case "error", "err":
		lvl = level.AllowError()
	}
	return lvl
}
