package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string) (*Logger, error) {
	loglevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	options := &slog.HandlerOptions{
		Level: loglevel,
	}

	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, options)),
	}, nil
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	if err != nil {
		return slog.LevelInfo, fmt.Errorf("invalid log level %q: %w", s, err)
	}
	return level, nil
}
