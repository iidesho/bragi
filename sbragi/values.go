package sbragi

import "golang.org/x/exp/slog"

type RedactedString string

func (RedactedString) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

func (r RedactedString) String() string {
	return string(r)
}

//This does not work. I want to make this work, but that has to wait. The issue is when it is a part of a struct and that struct gets logged. Creates a fake sense of security
