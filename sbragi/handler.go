package sbragi

/*
import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/iidesho/bragi"
	"github.com/iidesho/bragi/sbragi/mergedcontext"
)

type handler struct {
	slog.Handler
	ctx    context.Context
	cancel context.CancelFunc
}

func NewHandlerInFolder() (h handler, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	h = handler{
		ctx:    ctx,
		cancel: cancel,
	}
	handlerOpt := slog.HandlerOptions{
		AddSource: false,
		// Set a custom level to show all log output. The default value is
		// LevelInfo, which would drop Debug and Trace logs.
		Level: LevelInfo,

		ReplaceAttr: ReplaceAttr,
	}
	return
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.level <= level
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
*/
