package logger

import (
	"log/slog"
	"os"
)

func Setup() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     getLogLevel(),
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func getLogLevel() slog.Level {
	logLevelMap := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}

	envLevel := os.Getenv("LOG_LEVEL")
	if level, exists := logLevelMap[envLevel]; exists {
		return level
	}

	slog.Warn("Invalid LOG_LEVEL, defaulting to info", slog.String("provided", envLevel))
	return slog.LevelInfo
}
