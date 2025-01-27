package core

import (
	"log/slog"
)

// RedactString wraps a string that will be redacted in logs
type RedactString string

func (s RedactString) LogValue() slog.Value {
	str := string(s)
	if str == "" {
		return slog.StringValue("")
	}

	return slog.StringValue("[REDACTED]")
}
