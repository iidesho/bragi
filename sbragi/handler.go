package sbragi

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/cantara/bragi"
	"github.com/cantara/bragi/sbragi/mergedcontext"
)

const (
	LevelTrace   = slog.Level(-8)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelNotice  = slog.Level(2)
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

type handler struct {
	folder     string
	folderJson string
	fileHuman  *os.File
	fileJson   *os.File
	human      slog.Handler
	json       slog.Handler
	level      slog.Level
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewHandlerInFolder(path string) (h handler, err error) {
	path = strings.TrimSuffix(path, "/")
	ctx, cancel := context.WithCancel(context.Background())
	h = handler{
		folder:     path,
		folderJson: path + "/json",
		ctx:        ctx,
		cancel:     cancel,
	}
	if !bragi.FileExists(h.folder) {
		err = os.MkdirAll(h.folder, 0755)
		if err != nil {
			return
		}
	}
	if !bragi.FileExists(h.folderJson) {
		err = os.MkdirAll(h.folderJson, 0755)
		if err != nil {
			return
		}
	}
	h.fileHuman, h.fileJson, err = bragi.NewLogFiles(h.folder, h.folderJson)
	if err != nil {
		//bragi.AddError(err).Error("unable to create new logfiles")
		return
	}
	handlerOpt := slog.HandlerOptions{
		AddSource: false,
		// Set a custom level to show all log output. The default value is
		// LevelInfo, which would drop Debug and Trace logs.
		Level: LevelInfo,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			// Customize the name of the level key and the output string, including
			// custom level values.
			if a.Key == slog.LevelKey {
				// Handle custom level values.
				level := a.Value.Any().(slog.Level)

				// This could also look up the name from a map or other structure, but
				// this demonstrates using a switch statement to rename levels. For
				// maximum performance, the string values should be constants, but this
				// example uses the raw strings for readability.
				switch {
				case level < LevelDebug:
					a.Value = slog.StringValue("TRACE")
				case level < LevelInfo:
					a.Value = slog.StringValue("DEBUG")
				case level < LevelNotice:
					a.Value = slog.StringValue("INFO")
				case level < LevelWarning:
					a.Value = slog.StringValue("NOTICE")
				case level < LevelError:
					a.Value = slog.StringValue("WARNING")
				case level < LevelFatal:
					a.Value = slog.StringValue("ERROR")
				default:
					a.Value = slog.StringValue("FATAL")
				}
			}

			return a
		},
	}
	jsonHandleOpt := handlerOpt
	jsonHandleOpt.AddSource = true
	h.human = slog.NewTextHandler(h.fileHuman, &handlerOpt)
	h.json = slog.NewJSONHandler(h.fileJson, &jsonHandleOpt)
	go func() {
		nextDay := time.Now().UTC().AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 1, time.UTC)
		nextDayIn := nextDay.Sub(time.Now().UTC())
		rotateTicker := time.Tick(time.Second)
		rotateDayTicker := time.NewTicker(nextDayIn)
		truncateTaleTicker := time.Tick(time.Second * 5)
		firstDay := true
		slog.Info(fmt.Sprintf("all tickers for logger is created, next day is in: %v", nextDayIn), "next_day_in", nextDayIn)
		for {
			select {
			case <-ctx.Done():
				slog.Info("logger done selected. exiting")
				h.Cancel()
				return
			case <-rotateTicker:
				//Debug("logger rotate ticker selected")
				jsonStat, err := h.fileJson.Stat()
				if err != nil {
					slog.Log(ctx, LevelFatal, "unable to get json log file stats for rotation", "error", err.Error())
					continue
				}
				if jsonStat.Size() < 24*bragi.MB {
					//Debug("skipping rotate as filesize is less than 24KB base2. size is: ", jsonStat.Size(), " < ", 24*MB)
					continue
				}
				h.fileHuman, h.fileJson, err = bragi.Rotate(path, h.folderJson)
				if err != nil {
					slog.Log(ctx, LevelFatal, "unable to rotate", "error", err.Error())
					continue
				}
				h.human = slog.NewTextHandler(h.fileHuman, &handlerOpt)
				h.json = slog.NewJSONHandler(h.fileJson, &jsonHandleOpt)
			case <-rotateDayTicker.C:
				Debug("logger daily rotate ticker selected")
				if firstDay {
					firstDay = false
					rotateDayTicker.Reset(24 * time.Hour)
				}
				h.fileHuman, h.fileJson, err = bragi.Rotate(path, h.folderJson)
				if err != nil {
					slog.Log(ctx, LevelFatal, "unable to rotate", "error", err.Error())
					continue
				}
				h.human = slog.NewTextHandler(h.fileHuman, &handlerOpt)
				h.json = slog.NewJSONHandler(h.fileJson, &jsonHandleOpt)
			case <-truncateTaleTicker:
				Debug("logger truncate ticker selected")
				bragi.TruncateTale(h.folder)
				bragi.TruncateTale(h.folderJson)
			}
		}
	}()
	return
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	if h.level > level {
		return false
	}
	return true
}

func (h *handler) Handle(ctx context.Context, r slog.Record) (err error) {
	ctx, _ = mergedcontext.MergeContexts(h.ctx, ctx)
	err = h.human.Handle(ctx, r)
	if err != nil {
		return
	}
	return h.json.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.human.WithAttrs(attrs)
	h.json.WithAttrs(attrs)
	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &handler{
		folder:     h.folder,
		folderJson: h.folderJson,
		fileHuman:  h.fileHuman,
		fileJson:   h.fileJson,
		human:      h.human.WithGroup(name),
		json:       h.json.WithGroup(name),
		level:      h.level,
		ctx:        h.ctx,
		cancel:     h.cancel,
	}
}

func (h *handler) Cancel() {
	h.fileHuman.Close()
	h.fileJson.Close()
	h.cancel()
}

func (h *handler) MakeDefault() {
	slog.SetDefault(slog.New(h))
}

func (h *handler) SetLevel(level slog.Level) {
	h.level = level
}
