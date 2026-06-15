package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func setupOutputs(cfg *Config) (io.Writer, *os.File) {
	var writers []io.Writer
	var logFile *os.File
	for _, output := range cfg.Outputs {
		switch strings.ToLower(strings.TrimSpace(output)) {
		case "console":
			writers = append(writers, os.Stdout)
		case "file":
			logFile, err := setupOutputFile(cfg.Folder)
			if err != nil {
				panic(fmt.Errorf("failed to setup log file: %w", err))
			}
			writers = append(writers, logFile)
		default:
			fmt.Printf("warning: unknown logger output type: %s\n", output)
		}
	}
	if len(writers) == 0 {
		writers = append(writers, io.Discard)
	}
	return io.MultiWriter(writers...), logFile
}

func setupOutputFile(folder string) (*os.File, error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", folder, err)
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		folder,
		fmt.Sprintf("%s.log", timestamp),
	)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	return logFile, nil
}

func setupLevel(level string) slog.Level {
	var l slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		l = slog.LevelDebug
	case "WARN":
		l = slog.LevelWarn
	case "ERROR":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return l
}

func setupHandlerOpt(lvl slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     lvl,
		AddSource: lvl == slog.LevelDebug,
	}
}

func setupHandler(w io.Writer, opts *slog.HandlerOptions, format string) slog.Handler {
	switch strings.ToLower(format) {
	case "json":
		return slog.NewJSONHandler(w, opts)
	default:
		return slog.NewTextHandler(w, opts)
	}
}
