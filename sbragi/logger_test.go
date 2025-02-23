package sbragi_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/iidesho/bragi/sbragi"
)

var log = sbragi.WithLocalScope(sbragi.LevelDebug)

func TestLogger(t *testing.T) {
	h, err := sbragi.NewHandlerInFolder("./log")
	if err != nil {
		t.Error(err)
		return
	}
	h.SetLevel(sbragi.LevelTrace)
	defer h.Cancel()
	log, err := sbragi.NewLogger(&h)
	if err != nil {
		t.Error(err)
		return
	}
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.Trace("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.Debug("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 3")).Info("test")
	log.WithError(fmt.Errorf("simple error 3")).Info("test")
	log.Info("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.Notice("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.Warning("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 6")).Error("test")
	log.WithError(fmt.Errorf("simple error 6")).Error("test")
	log.Error("test")
	/*
		log.WithError(fmt.Errorf("simple error 7")).Fatal("test")
		log.Fatal("test")
	*/
}

func TestPackageLevelLogger(t *testing.T) {
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.Trace("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.Debug("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 3")).Info("test")
	log.WithError(fmt.Errorf("simple error 3")).Info("test")
	log.Info("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.Notice("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.Warning("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 6")).Error("test")
	log.WithError(fmt.Errorf("simple error 6")).Error("test")
	log.Error("test")
	/*
		log.WithError(fmt.Errorf("simple error 7")).Fatal("test")
		log.Fatal("test")
	*/
}

func TestPackageFunctionLevelLogger(t *testing.T) {
	log := sbragi.GetDefaultLogger().WithLocalScope(sbragi.LevelNotice)
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.WithError(fmt.Errorf("simple error 1")).Trace("test")
	log.Trace("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.WithError(fmt.Errorf("simple error 2")).Debug("test")
	log.Debug("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 3")).Info("test")
	log.WithError(fmt.Errorf("simple error 3")).Info("test")
	log.Info("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.WithError(fmt.Errorf("simple error 4")).Notice("test")
	log.Notice("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.WithError(fmt.Errorf("simple error 5")).Warning("test")
	log.Warning("test")
	log.WithoutEscalation().WithError(fmt.Errorf("simple error 6")).Error("test")
	log.WithError(fmt.Errorf("simple error 6")).Error("test")
	log.Error("test")
	/*
		log.WithError(fmt.Errorf("simple error 7")).Fatal("test")
		log.Fatal("test")
	*/
}

func TestDebugLogger(t *testing.T) {
	dl, err := sbragi.NewDebugLogger()
	if err != nil {
		t.Error(err)
		return
	}
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 1")).Trace("test")
	dl.WithError(fmt.Errorf("simple error 1")).Trace("test")
	dl.Trace("test")
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 2")).Debug("test")
	dl.WithError(fmt.Errorf("simple error 2")).Debug("test")
	dl.Debug("test")
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 3")).Info("test")
	dl.WithError(fmt.Errorf("simple error 3")).Info("test")
	dl.Info("test")
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 4")).Notice("test")
	dl.WithError(fmt.Errorf("simple error 4")).Notice("test")
	dl.Notice("test")
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 5")).Warning("test")
	dl.WithError(fmt.Errorf("simple error 5")).Warning("test")
	dl.Warning("test")
	dl.WithoutEscalation().WithError(fmt.Errorf("simple error 6")).Error("test")
	dl.WithError(fmt.Errorf("simple error 6")).Error("test")
	dl.Error("test")
	/*
		log.WithError(fmt.Errorf("simple error 7")).Fatal("test")
		log.Fatal("test")
	*/
}

func TestDynamicLogLevel(t *testing.T) {
	log := sbragi.WithLocalScope(sbragi.LevelInfo)
	f, err := os.CreateTemp(".", "scope_levels-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sbragi.AttachDynamicScopes(f.Name())
	log.Info("attacked dynamic scope", "path", f.Name())
	log.Debug("testing before file change")
	f.WriteString("github.com/iidesho/bragi: debug\n")
	for i := range 15 {
		log.Debug("file changed, waiting for change to take effect", "i", i)
		time.Sleep(time.Microsecond * 10)
	}
	log.Info("15 debug messages printed")
	f.WriteString("github.com/iidesho/bragi/sbragi_test: error\n")
	for i := range 10 {
		log.Debug("file changed, waiting for change to take effect", "i", i)
		time.Sleep(time.Microsecond * 10)
	}
	log.Error("10 debug messages printed")
}

func TestDefaultDebugLogger(t *testing.T) {
	dl, err := sbragi.NewDebugLogger()
	if err != nil {
		t.Error(err)
		return
	}
	dl.SetDefault()
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 1")).Trace("test")
	sbragi.WithError(fmt.Errorf("simple error 1")).Trace("test")
	sbragi.Trace("test")
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 2")).Debug("test")
	sbragi.WithError(fmt.Errorf("simple error 2")).Debug("test")
	sbragi.Debug("test")
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 3")).Info("test")
	sbragi.WithError(fmt.Errorf("simple error 3")).Info("test")
	sbragi.Info("test")
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 4")).Notice("test")
	sbragi.WithError(fmt.Errorf("simple error 4")).Notice("test")
	sbragi.Notice("test")
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 5")).Warning("test")
	sbragi.WithError(fmt.Errorf("simple error 5")).Warning("test")
	sbragi.Warning("test")
	sbragi.WithoutEscalation().WithError(fmt.Errorf("simple error 6")).Error("test")
	sbragi.WithError(fmt.Errorf("simple error 6")).Error("test")
	sbragi.Error("test")
	/*
		log.WithError(fmt.Errorf("simple error 7")).Fatal("test")
		log.Fatal("test")
	*/
}

func BenchmarkLogger(b *testing.B) {
	log, err := sbragi.NewLogger(slog.NewJSONHandler(os.Stdout, nil))
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		log.Error("bench", "number", i)
	}
}

func BenchmarkLoggerWHandler(b *testing.B) {
	h, err := sbragi.NewHandlerInFolder("./log")
	if err != nil {
		b.Error(err)
		return
	}
	defer h.Cancel()
	log, err := sbragi.NewLogger(&h)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		log.Error("bench", "number", i)
	}
}
