package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger with the specified configuration
func Init(level string, format string) error {
	var config zap.Config

	// Set log level
	logLevel, err := parseLogLevel(level)
	if err != nil {
		logLevel = zap.InfoLevel
	}

	// Configure based on format
	if format == "json" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Build logger
	var err2 error
	log, err2 = config.Build()
	if err2 != nil {
		return err2
	}

	// Replace global logger
	zap.ReplaceGlobals(log)

	return nil
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		// Initialize with defaults if not initialized
		if err := Init("info", "text"); err != nil {
			// Fallback to a basic logger if initialization fails
			log, _ = zap.NewDevelopment()
		}
	}
	return log
}

// parseLogLevel converts string to zap log level
func parseLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "fatal":
		return zap.FatalLevel, nil
	case "panic":
		return zap.PanicLevel, nil
	default:
		return zap.InfoLevel, nil
	}
}

// Convenience functions for common logging operations

// Info logs info level message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs error level message
func Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	GetLogger().Error(msg, fields...)
}

// Debug logs debug level message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn logs warning level message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal logs fatal level message and exits
func Fatal(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	GetLogger().Fatal(msg, fields...)
}

// WithContext creates a logger with request context fields
func WithContext(requestID, userID, method, path string) *zap.Logger {
	return GetLogger().With(
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.String("method", method),
		zap.String("path", path),
	)
}

// WithNotification creates a logger with notification-specific fields
func WithNotification(notificationID, recipient, notificationType string) *zap.Logger {
	return GetLogger().With(
		zap.String("notification_id", notificationID),
		zap.String("recipient", recipient),
		zap.String("notification_type", notificationType),
	)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}
