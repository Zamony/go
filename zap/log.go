package zap

import (
	"log/slog"
	"runtime"
	"time"
)

type Logger struct {
	Logger *slog.Logger
}

func (l *Logger) log(level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.Logger.Enabled(nil, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log, LevelCall]

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.AddAttrs(attrs...)
	_ = l.Logger.Handler().Handle(nil, record)
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.log(slog.LevelDebug, msg, attrs...)
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.log(slog.LevelInfo, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	l.log(slog.LevelWarn, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.log(slog.LevelError, msg, attrs...)
}

func (l *Logger) With(attrs ...slog.Attr) Logger {
	if len(attrs) == 0 {
		return *l
	}

	return Logger{Logger: slog.New(l.Logger.Handler().WithAttrs(attrs))}
}

var ErrorStackMarshaler func(error) string

func Error(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "nil")
	}

	attr := slog.String("error", err.Error())
	if ErrorStackMarshaler == nil {
		return attr
	}

	stack := ErrorStackMarshaler(err)
	if stack == "" {
		return attr
	}

	return slog.Group("", attr, slog.String("errstack", stack))
}

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func Time(key string, value time.Time) slog.Attr {
	return slog.Time(key, value)
}

func Uint64(key string, value uint64) slog.Attr {
	return slog.Uint64(key, value)
}
