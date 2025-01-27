package sbragi

import (
	"fmt"
	"log/slog"
	"strings"
)

var log = WithLocalScope(LevelWarning)

const (
	LevelTrace   = slog.Level(-8)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelNotice  = slog.Level(2)
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
	LevelUnknown = slog.Level(13)
)

func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Handle custom level values.
		// fmt.Println(a.Value)
		// level := a.Value.Any().(slog.Level)
		// a.Value = slog.StringValue(LevelToString(level))
		a.Value = slog.StringValue(ppSLogString(a.Value.String()))
	}

	return a
}

func LevelToString(level slog.Level) string {
	switch level {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelNotice:
		return "NOTICE"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelUnknown:
		return "UNKNOWN"
	default:
		return level.String()
	}
}

func ppSLogString(level string) string {
	switch level {
	case LevelTrace.String():
		return "TRACE"
	case LevelDebug.String():
		return "DEBUG"
	case LevelInfo.String():
		return "INFO"
	case LevelNotice.String():
		return "NOTICE"
	case LevelWarning.String():
		return "WARNING"
	case LevelError.String():
		return "ERROR"
	case LevelFatal.String():
		return "FATAL"
	case LevelUnknown.String():
		return "UNKNOWN"
	default:
		return level
	}
}

func StringToLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return LevelTrace
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "NOTICE":
		return LevelNotice
	case "WARNING":
		return LevelWarning
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		fmt.Println("unknown", level, strings.ToUpper(level), "ERROR")
		return LevelUnknown
	}
}
