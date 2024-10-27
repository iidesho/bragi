package sbragi

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/iidesho/bragi"
	"github.com/iidesho/bragi/sbragi/mergedcontext"
)

type fileHandler struct {
	human      slog.Handler
	json       slog.Handler
	ctx        context.Context
	fileHuman  *os.File
	fileJson   *os.File
	cancel     context.CancelFunc
	folder     string
	folderJson string
	level      slog.Level
}

func NewHandlerInFolder(path string) (h fileHandler, err error) {
	path = strings.TrimSuffix(path, "/")
	ctx, cancel := context.WithCancel(context.Background())
	h = fileHandler{
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

		ReplaceAttr: ReplaceAttr,
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
		slog.Info(
			fmt.Sprintf("all tickers for logger is created, next day is in: %v", nextDayIn),
			"next_day_in",
			nextDayIn,
		)
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
					slog.Log(
						ctx,
						LevelFatal,
						"unable to get json log file stats for rotation",
						"error",
						err.Error(),
					)
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

func (h *fileHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.level <= level
}

func (h *fileHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	ctx, _ = mergedcontext.MergeContexts(h.ctx, ctx)
	err = h.human.Handle(ctx, r)
	if err != nil {
		return
	}
	return h.json.Handle(ctx, r)
}

func (h *fileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.human.WithAttrs(attrs)
	h.json.WithAttrs(attrs)
	return h
}

func (h *fileHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &fileHandler{
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

func (h *fileHandler) Cancel() {
	h.fileHuman.Close()
	h.fileJson.Close()
	h.cancel()
}

func (h *fileHandler) MakeDefault() {
	slog.SetDefault(slog.New(h))
}

func (h *fileHandler) SetLevel(level slog.Level) {
	h.level = level
}
