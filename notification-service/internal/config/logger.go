package config

import (
	"log/slog"
	"os"
	"strings"
)

func ParseLog(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func InitLogger() *slog.Logger {
	logLevelENV := os.Getenv("LOG_LEVEL")
	if logLevelENV == "" {
		logLevelENV = "info"
	}

	logLevel := ParseLog(logLevelENV)
	handlersLogger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	logger := slog.New(handlersLogger)

	slog.Info("logger инициализирован: ", "level", logLevel.String())

	return logger
}
