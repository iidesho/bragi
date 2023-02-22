package sbragi

import (
	"golang.org/x/exp/slog"
	"os"
)

var defaultLogger, _ = NewLogger(slog.NewTextHandler(os.Stdout))

type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Notice(msg string, args ...any)
	Warning(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	WithError(err error) Logger
	SetDefault()
}

type logger struct {
	handler slog.Handler
	log     *slog.Logger
	depth   int
	err     error
}

func NewLogger(handler slog.Handler) (Logger, error) {
	return newLogger(2, handler)
}

func newLogger(depth int, handler slog.Handler) (Logger, error) {
	return logger{
		handler: handler,
		log:     slog.New(handler),
		depth:   depth,
	}, nil
}

func (l logger) SetDefault() {
	defaultLogger = l
}

func (l logger) Trace(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelTrace) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelTrace, msg, args...)
}

func (l logger) Debug(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelDebug) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelDebug, msg, args...)
}

func (l logger) Info(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelInfo) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelInfo, msg, args...)
}

func (l logger) Notice(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelNotice) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelNotice, msg, args...)
}

func (l logger) Warning(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelWarning) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelWarning, msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelError) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelError, msg, args...)
}

func (l logger) Fatal(msg string, args ...any) {
	if !l.handler.Enabled(nil, LevelFatal) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	l.log.LogDepth(l.depth, LevelFatal, msg, args...)
	panic(msg)
}

func (l logger) WithError(err error) Logger {
	l.err = err
	return l
}

func Trace(msg string, args ...any) {
	defaultLogger.Trace(msg, args...)
}
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}
func Notice(msg string, args ...any) {
	defaultLogger.Notice(msg, args...)
}
func Warning(msg string, args ...any) {
	defaultLogger.Warning(msg, args...)
}
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}
func Fatal(msg string, args ...any) {
	defaultLogger.Fatal(msg, args...)
}
func WithError(err error) Logger {
	return defaultLogger.WithError(err)
}
