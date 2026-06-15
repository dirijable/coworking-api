package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

//var (
//	mu sync.Mutex
//)

func NewLogger(log *slog.Logger, file *os.File) *Logger {
	return &Logger{
		Logger: log,
		file:   file,
	}
}

func NewFromConfig(cfg Config) *Logger {
	writer, logFile := setupOutputs(&cfg)
	slogLvl := setupLevel(cfg.Level)
	handlerOpt := setupHandlerOpt(slogLvl)
	handler := setupHandler(writer, handlerOpt, cfg.Format)
	logger := slog.New(handler)
	return NewLogger(logger, logFile)
}

func (l *Logger) Close() error {
	//mu.Lock()
	//defer mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
