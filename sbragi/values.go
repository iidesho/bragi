package sbragi

import "golang.org/x/exp/slog"

type RedactedString string

func (RedactedString) LogValue() slog.Value {
	return slog.StringValue("REDACTED")
}

func (r RedactedString) String() string {
	return string(r)
}
