package stlog

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
