package log

import (
	gofmt "fmt"
	"math"
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
	unixTimeFormat = float64(time.Now().UTC().UnixNano()) / math.Pow(10, 9)
	defaultLogger  = New(&Config{})
)

type Logger struct {
	debug log.Logger
	info  log.Logger
	warn  log.Logger
	err   log.Logger

	lvl    string
	format string
}

func (l *Logger) Debug(fmt string, v ...interface{}) {
	err := l.debug.Log(prepareLog(fmt, v...))
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Info(fmt string, v ...interface{}) {
	err := l.info.Log(prepareLog(fmt, v...))
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Warn(fmt string, v ...interface{}) {
	err := l.warn.Log(prepareLog(fmt, v...))
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Error(fmt string, v ...interface{}) {
	err := l.err.Log(prepareLog(fmt, v...))
	if err != nil {
		panic(err)
	}
}

func (l *Logger) SetLevel(lvl string) {
	*l = *(New(&Config{Level: lvl, Format: l.format}))
}

func (l *Logger) SetFormat(format string) {
	*l = *(New(&Config{Level: l.lvl, Format: format}))
}

func Debug(fmt string, v ...interface{}) {
	defaultLogger.Debug(fmt, v...)
}

func Info(fmt string, v ...interface{}) {
	defaultLogger.Info(fmt, v...)
}

func Warn(fmt string, v ...interface{}) {
	defaultLogger.Warn(fmt, v...)
}

func Error(fmt string, v ...interface{}) {
	defaultLogger.Error(fmt, v...)
}

func SetLevel(lvl string) {
	defaultLogger.SetLevel(lvl)
}

func SetFormat(format string) {
	defaultLogger.SetFormat(format)
}

type Config struct {
	Level  string
	Format string
}

func New(cfg *Config) *Logger {
	var l log.Logger
	if cfg.Format == "json" {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}

	l = log.With(level.NewFilter(l, parseLevel(cfg.Level)), "ts", unixTimeFormat, "date", timestampFormat, "logger", "metrics", "caller", log.DefaultCaller)

	return &Logger{debug: level.Debug(l), info: level.Info(l), warn: level.Warn(l), err: level.Error(l), lvl: cfg.Level, format: cfg.Format}
}

func prepareLog(fmt string, v ...interface{}) (string, string) {
	return "message", gofmt.Sprintf(fmt, v...)
}

func parseLevel(l string) level.Option {
	lvl := level.AllowDebug()

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
