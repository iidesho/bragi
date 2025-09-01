package sbragi

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	contextkeys "github.com/iidesho/gober/contextKeys"
	"go.opentelemetry.io/otel/trace"
)

// Using debug logger as default logger as we use scopes for granulaity
// If we could dynamically set level based on config at runtime, then we could
// Use propper logging level info as default
var defaultLogger, _ = NewDebugLogger()

func GetDefaultLogger() DefaultLogger {
	return &defaultLogger
}

type DefaultLogger interface {
	ContextLogger
	WithLocalScope(defaultLevel slog.Level) ContextLogger
}

type ContextLogger interface {
	ErrorLogger
	WithContext(ctx context.Context) ErrorLogger
}

type ErrorLogger interface {
	BaseLogger
	WithError(err error) BaseLogger
	WithErrorFunc(errf func() error) BaseLogger
	WithoutEscalation() ErrorLogger
}

type BaseLogger interface {
	Trace(msg string, args ...any) bool
	Debug(msg string, args ...any) bool
	Printf(format string, args ...any)
	Info(msg string, args ...any) bool
	Notice(msg string, args ...any) bool
	Warning(msg string, args ...any) bool
	Error(msg string, args ...any) bool
	Level(lvl slog.Level, msg string, args ...any) bool
	Fatal(msg string, args ...any)
}

type scopeLevel struct {
	scope string
	level slog.Level
}

type logger struct {
	handler   slog.Handler
	ctx       context.Context
	err       error
	scopes    *[]scopeLevel
	scopesMut *sync.RWMutex
	slog      *slog.Logger
	errf      func() error
	scope     string
	level     slog.Level
	depth     int
	escalate  bool
	withError bool
}

func NewLogger(handler slog.Handler) (logger, error) {
	return newLogger(handler)
}

func NewDebugLogger() (logger, error) {
	return NewLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       LevelDebug,
		ReplaceAttr: ReplaceAttr,
	}))
}

func NewTraceLogger() (logger, error) {
	return NewLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       LevelTrace,
		ReplaceAttr: ReplaceAttr,
	}))
}

func newLogger(handler slog.Handler) (logger, error) {
	var scopes []scopeLevel
	return logger{
		handler:   handler,
		slog:      slog.New(handler),
		ctx:       context.Background(), // This is just a temporaty context
		escalate:  true,
		scopes:    &scopes,
		scopesMut: &sync.RWMutex{},
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

func (l logger) Printf(format string, args ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, args...))
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

func (l logger) Level(lvl slog.Level, msg string, args ...any) bool {
	return l.log(lvl, msg, args...)
}

func (l logger) Fatal(msg string, args ...any) {
	if l.log(LevelFatal, msg, args...) || !l.withError {
		panic(fmt.Sprint(msg, args))
	}
}

func (l logger) WithContext(ctx context.Context) ErrorLogger {
	l.ctx = ctx
	return l
}

func (l logger) WithError(err error) BaseLogger {
	l.err = err
	l.withError = true
	/*
		if l.scope == "" {
			l.depth--
		}
	*/
	// l.depth--
	return l
}

func (l logger) WithErrorFunc(errf func() error) BaseLogger {
	l.errf = errf
	l.withError = true
	/*
		if l.scope == "" {
			l.depth--
		}
	*/
	// l.depth--
	return l
}

func (l logger) WithoutEscalation() ErrorLogger {
	l.escalate = false
	return l
}

func (l logger) WithLocalScope(defaultLevel slog.Level) ContextLogger {
	return l.withLocalScope(defaultLevel)
}

func (l logger) withLocalScope(defaultLevel slog.Level) ContextLogger {
	pc, _, _, ok := runtime.Caller(2) // This is a super ugly hack :/
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		Fatal("could not get runtime information about caller")
	}
	l.scope = strings.TrimSuffix(details.Name(), ".init")
	l.level = defaultLevel
	l.Trace("local scope added", "level", LevelToString(defaultLevel), "scope", l.scope)
	//l.depth++
	/*
		frames := runtime.CallersFrames([]uintptr{pc})

		// Loop to get frames.
		// A fixed number of PCs can expand to an indefinite number of Frames.
		frame, _ := frames.Next() //Ignoring more as we only care about caller
		fuctionParts := strings.Split(frame.Function, ".")
		fmt.Printf(
			"function %s, package %s\n",
			frame.Function,
			strings.Join(fuctionParts[:len(fuctionParts)-1], "."),
		)
	*/
	return l
}

func readScopeConfig(f io.ReadCloser) []scopeLevel {
	defer log.WithErrorFunc(f.Close).Trace("closed scopes config file")
	bf := bufio.NewScanner(f)
	scopes := []scopeLevel{}
	for bf.Scan() {
		t := bf.Text()
		data := strings.SplitN(t, ":", 2)
		l := strings.TrimSpace(data[1])
		scopes = append(scopes, scopeLevel{
			scope: strings.TrimSpace(data[0]),
			level: StringToLevel(l),
		})
	}
	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i].scope > scopes[j].scope
	})
	return scopes
}

func AttachDynamicScopes(config string) {
	log := &defaultLogger
	if *log.scopes != nil {
		log.Error("trying to attach new dymanic configuration", "config", config)
		return
	}

	f, err := os.OpenFile(config, os.O_CREATE|os.O_RDONLY, 0640)
	if log.WithError(err).Error("could not open scopes config file") {
		return
	}
	log.scopesMut.Lock()
	*log.scopes = readScopeConfig(f)
	log.scopesMut.Unlock()

	watcher, err := fsnotify.NewWatcher()
	if log.WithError(err).Error("could not create watcher") {
		return
	}
	if log.WithError(watcher.AddWith(f.Name())).Error("adding watcher") {
		return
	}
	go func(events chan fsnotify.Event) {
		for e := range events {
			if !e.Op.Has(fsnotify.Write) {
				continue
			}
			log.Info("Config file changed, reloading...")
			f, err := os.OpenFile(config, os.O_RDONLY, 0640)
			if log.WithError(err).Error("could not open scopes config file") {
				continue
			}
			log.scopesMut.Lock()
			*log.scopes = readScopeConfig(f)
			log.scopesMut.Unlock()
		}
	}(watcher.Events)
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

func Level(lvl slog.Level, msg string, args ...any) bool {
	return defaultLogger.Level(lvl, msg, args...)
}

func WithError(err error) BaseLogger {
	l := defaultLogger
	l.depth--
	return l.WithError(err)
}

func WithErrorFunc(errf func() error) BaseLogger {
	l := defaultLogger
	l.depth--
	return l.WithErrorFunc(errf)
}

func WithoutEscalation() ErrorLogger {
	l := defaultLogger
	l.depth--
	return l.WithoutEscalation()
}

func WithLocalScope(defaultLevel slog.Level) ContextLogger {
	l := defaultLogger
	// l.depth--
	return l.withLocalScope(defaultLevel)
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l logger) log(level slog.Level, msg string, args ...any) (hadError bool) {
	// loggedError = l.err != nil
	if l.errf != nil {
		l.err = l.errf()
	}
	if l.err != nil {
		hadError = true
	}
	if l.err != nil {
		if l.escalate &&
			level < LevelNotice { // This is not intuitive :/ a escalate from level might solve this
			level = LevelError
		}
	} else {
		if level >= LevelNotice && l.withError || l.escalate && l.withError {
			return // false // Return early if level is error and there was no error
		}
	}
	if !l.handler.Enabled(l.ctx, level) {
		return // false
	}
	spanCTX := trace.SpanContextFromContext(l.ctx)
	if spanCTX.IsValid() {
		args = append(
			args,
			"trace_id",
			spanCTX.TraceID().String(),
			"span_id",
			spanCTX.SpanID().String(),
		)
	}
	for _, ctxKey := range contextkeys.Keys {
		if tid := l.ctx.Value(ctxKey); tid != nil {
			args = append([]any{ctxKey.String(), tid}, args...)
		}
	}
	if l.scope != "" {
		// return ealy if the loggers local scope reauires a higher level than requested
		if *l.scopes != nil {
			l.scopesMut.RLock()
			for _, scope := range *l.scopes {
				if strings.HasPrefix(l.scope, scope.scope) {
					l.level = scope.level
					break
				}
			}
			l.scopesMut.RUnlock()
		}
		if l.level > level {
			return // false
		}
		args = append([]any{"scope", l.scope}, args...)
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
	return
}
