package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger *slog.Logger

// Config for logger
type Config struct {
	Level      string // debug, info, warn, error
	FilePath   string // path to log file
	MaxSize    int    // max size in MB before rotation
	MaxBackups int    // max number of old log files to keep
	MaxAge     int    // max days to retain old log files
	Compress   bool   // compress rotated files
}

// DefaultConfig returns default logger config
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

// Init initializes the global logger
func Init(cfg Config) error {
	// Create logs directory if not exists
	logDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Setup lumberjack for log rotation
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// Parse log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler with JSON format
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)

	return nil
}

// GetLogger returns the default logger
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		// Fallback to default slog if not initialized
		return slog.Default()
	}
	return defaultLogger
}

// Context key for request ID
type contextKey string

const RequestIDKey contextKey = "request_id"

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID gets request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// Helper functions with context support

// InfoContext logs info level with context
func InfoContext(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	GetLogger().Info(msg, args...)
}

// WarnContext logs warn level with context
func WarnContext(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	GetLogger().Warn(msg, args...)
}

// ErrorContext logs error level with context
func ErrorContext(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	GetLogger().Error(msg, args...)
}

// DebugContext logs debug level with context
func DebugContext(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	GetLogger().Debug(msg, args...)
}

// appendRequestID appends request_id to args if present in context
func appendRequestID(ctx context.Context, args []any) []any {
	if requestID := GetRequestID(ctx); requestID != "" {
		args = append(args, "request_id", requestID)
	}
	return args
}

// Info logs info level
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs warn level
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs error level
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// Debug logs debug level
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}
