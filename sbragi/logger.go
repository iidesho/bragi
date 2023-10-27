package sbragi

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var defaultLogger, _ = NewLogger(slog.NewTextHandler(os.Stdout, nil))

type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Notice(msg string, args ...any)
	Warning(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	WithError(err error) Logger
	WithoutEscalation() Logger
	SetDefault()
}

type logger struct {
	handler  slog.Handler
	slog     *slog.Logger
	depth    int
	ctx      context.Context
	escalate bool
	err      error
}

func NewLogger(handler slog.Handler) (logger, error) {
	return newLogger(handler)
}

func NewDebugLogger() (logger, error) {
	return NewLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
	}))
}

func newLogger(handler slog.Handler) (logger, error) {
	return logger{
		handler:  handler,
		slog:     slog.New(handler),
		ctx:      context.Background(),
		escalate: true,
	}, nil
}

func (l logger) SetDefault() {
	l.depth++
	defaultLogger = l
}

func (l logger) Trace(msg string, args ...any) {
	l.log(LevelTrace, msg, args...)
}

func (l logger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

func (l logger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

func (l logger) Notice(msg string, args ...any) {
	l.log(LevelNotice, msg, args...)
}

func (l logger) Warning(msg string, args ...any) {
	l.log(LevelWarning, msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}

func (l logger) Fatal(msg string, args ...any) {
	l.log(LevelFatal, msg, args...)
	panic(fmt.Sprint(msg, args))
}

func (l logger) WithError(err error) Logger {
	l.err = err
	//l.depth--
	return l
}

func (l logger) WithoutEscalation() Logger {
	l.escalate = false
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
	l := defaultLogger
	l.depth--
	return l.WithError(err)
}
func WithoutEscalation() Logger {
	l := defaultLogger
	l.depth--
	return l.WithoutEscalation()
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l logger) log(level slog.Level, msg string, args ...any) {
	if l.escalate && l.err != nil && level < LevelError {
		level = LevelError
	}
	if !l.handler.Enabled(l.ctx, level) {
		return
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3+l.depth, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.handler.Handle(l.ctx, r)
}
