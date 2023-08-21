package sbragi

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	h, err := NewHandlerInFolder("./log")
	if err != nil {
		t.Error(err)
		return
	}
	h.SetLevel(LevelTrace)
	defer h.Cancel()
	log, err := NewLogger(&h)
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

func BenchmarkLogger(b *testing.B) {
	log, err := newLogger(slog.NewJSONHandler(os.Stdout, nil))
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		log.Error("bench", "number", i)
	}
}

func BenchmarkLoggerWHandler(b *testing.B) {
	h, err := NewHandlerInFolder("./log")
	if err != nil {
		b.Error(err)
		return
	}
	defer h.Cancel()
	log, err := newLogger(&h)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		log.Error("bench", "number", i)
	}
}
