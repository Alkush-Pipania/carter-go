package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger instance
var Log *zap.Logger

// Sugar is the sugared logger for printf-style logging
var Sugar *zap.SugaredLogger

// Context keys for contextual logging
type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
	UserIDKey    ctxKey = "user_id"
)

// Config holds logger configuration
type Config struct {
	Env   string // "development" or "production"
	Level string // "debug", "info", "warn", "error"
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	var zapConfig zap.Config

	if cfg.Env == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	}

	// Set log level from config
	zapConfig.Level = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

	// Build the logger
	var err error
	Log, err = zapConfig.Build(
		zap.AddCallerSkip(1), // Skip wrapper functions in stack trace
	)
	if err != nil {
		// Fallback to basic logger if config fails
		Log, _ = zap.NewProduction()
	}

	Sugar = Log.Sugar()
}

// parseLevel converts string level to zapcore.Level
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel // Default to info
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		Log.Sync()
	}
}

// WithContext returns a logger with contextual fields from context
func WithContext(ctx context.Context) *zap.Logger {
	logger := Log

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With(zap.String("request_id", requestID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		logger = logger.With(zap.String("user_id", userID))
	}

	return logger
}

// WithFields returns a logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}

// --- Convenience functions for structured logging ---

// Info logs an info message with optional fields
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Debug logs a debug message with optional fields
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Warn logs a warning message with optional fields
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error logs an error message with optional fields
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
	os.Exit(1)
}
