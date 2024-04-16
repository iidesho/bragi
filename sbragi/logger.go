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

func DefaultLogger() Logger {
	return defaultLogger
}

type Logger interface {
	ErrorLogger
	WithError(err error) ErrorLogger
	WithErrorFunc(errf func() error) ErrorLogger
	WithoutEscalation() Logger
	SetDefault()
}

type ErrorLogger interface {
	Trace(msg string, args ...any) bool
	Debug(msg string, args ...any) bool
	Info(msg string, args ...any) bool
	Notice(msg string, args ...any) bool
	Warning(msg string, args ...any) bool
	Error(msg string, args ...any) bool
	Fatal(msg string, args ...any)
}

type logger struct {
	handler   slog.Handler
	slog      *slog.Logger
	depth     int
	ctx       context.Context
	escalate  bool
	err       error
	errf      func() error
	withError bool
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

func (l logger) Trace(msg string, args ...any) bool {
	return l.log(LevelTrace, msg, args...)
}

func (l logger) Debug(msg string, args ...any) bool {
	return l.log(LevelDebug, msg, args...)
}

func (l logger) Info(msg string, args ...any) bool {
	return l.log(LevelInfo, msg, args...)
}

func (l logger) Notice(msg string, args ...any) bool {
	return l.log(LevelNotice, msg, args...)
}

func (l logger) Warning(msg string, args ...any) bool {
	return l.log(LevelWarning, msg, args...)
}

func (l logger) Error(msg string, args ...any) bool {
	return l.log(LevelError, msg, args...)
}

func (l logger) Fatal(msg string, args ...any) {
	if l.log(LevelFatal, msg, args...) || !l.withError {
		panic(fmt.Sprint(msg, args))
	}
}

func (l logger) WithError(err error) ErrorLogger {
	l.err = err
	l.withError = true
	//l.depth--
	return l
}
func (l logger) WithErrorFunc(errf func() error) ErrorLogger {
	l.errf = errf
	l.withError = true
	//l.depth--
	return l
}

func (l logger) WithoutEscalation() Logger {
	l.escalate = false
	return l
}

func Trace(msg string, args ...any) bool {
	return defaultLogger.Trace(msg, args...)
}
func Debug(msg string, args ...any) bool {
	return defaultLogger.Debug(msg, args...)
}
func Info(msg string, args ...any) bool {
	return defaultLogger.Info(msg, args...)
}
func Notice(msg string, args ...any) bool {
	return defaultLogger.Notice(msg, args...)
}
func Warning(msg string, args ...any) bool {
	return defaultLogger.Warning(msg, args...)
}
func Error(msg string, args ...any) bool {
	return defaultLogger.Error(msg, args...)
}
func Fatal(msg string, args ...any) {
	defaultLogger.Fatal(msg, args...)
}
func WithError(err error) ErrorLogger {
	l := defaultLogger
	l.depth--
	return l.WithError(err)
}
func WithErrorFunc(errf func() error) ErrorLogger {
	l := defaultLogger
	l.depth--
	return l.WithErrorFunc(errf)
}
func WithoutEscalation() Logger {
	l := defaultLogger
	l.depth--
	return l.WithoutEscalation()
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l logger) log(level slog.Level, msg string, args ...any) (loggedError bool) {
	//loggedError = l.err != nil
	if l.errf != nil {
		l.err = l.errf()
	}
	if l.err != nil {
		if l.escalate && level < LevelError {
			level = LevelError
		}
	} else {
		if level >= LevelError && l.withError {
			return false //Return early if level is error and there was no error
		}
	}
	if !l.handler.Enabled(l.ctx, level) {
		return false
	}
	if l.err != nil {
		args = append([]any{"error", l.err}, args...)
		loggedError = true
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3+l.depth, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.handler.Handle(l.ctx, r)
	return
}
