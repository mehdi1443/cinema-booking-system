package logger

import (
	"log/slog"
	"strings"
)

type Config struct {
	ServiceName string
	Environment string

	Level string

	AddSource bool
}

func (c Config) slogLevel() slog.Level {
	switch strings.ToLower(strings.TrimSpace(c.Level)) {
	case "debug":
		return slog.LevelDebug

	case "warn", "warning":
		return slog.LevelWarn

	case "error":
		return slog.LevelError

	case "info", "":
		return slog.LevelInfo

	default:
		return slog.LevelInfo
	}
}
