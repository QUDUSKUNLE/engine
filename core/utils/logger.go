package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// LogConfig holds the configuration for the logger
type LogConfig struct {
	Level       string
	OutputPath  string
	Development bool
}

// defaultConfig provides default logging configuration
func defaultConfig() LogConfig {
	return LogConfig{
		Level:       "info",
		OutputPath:  "stdout",
		Development: false,
	}
}

// InitLogger initializes the zap logger with the given configuration
// If no config is provided, default configuration is used
func InitLogger(config ...LogConfig) error {
	cfg := defaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Configure logging level
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure output paths
	var outputPaths []string
	if cfg.OutputPath == "stdout" {
		outputPaths = []string{"stdout"}
	} else {
		outputPaths = []string{cfg.OutputPath, "stdout"}
	}

	// Create the logger configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.Development,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	// Create the logger
	var err error
	Logger, err = zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	// Create sugar logger for convenience
	Sugar = Logger.Sugar()

	// Sync logger on application shutdown
	go func() {
		c := make(chan os.Signal, 1)
		<-c
		Logger.Sync()
	}()

	Logger.Info("Logger initialized",
		zap.String("level", cfg.Level),
		zap.String("output", cfg.OutputPath),
		zap.Bool("development", cfg.Development),
		zap.Time("init_time", time.Now()),
	)

	return nil
}

// LogField represents a key-value pair for structured logging
type LogField struct {
	Key   string
	Value interface{}
}

// Debug logs a message at debug level
func Debug(msg string, fields ...LogField) {
	if len(fields) == 0 {
		Sugar.Debug(msg)
		return
	}
	zapFields := convertToZapFields(fields)
	Logger.Debug(msg, zapFields...)
}

// Info logs a message at info level
func Info(msg string, fields ...LogField) {
	if len(fields) == 0 {
		Sugar.Info(msg)
		return
	}
	zapFields := convertToZapFields(fields)
	Logger.Info(msg, zapFields...)
}

// Warn logs a message at warn level
func Warn(msg string, fields ...LogField) {
	if len(fields) == 0 {
		Sugar.Warn(msg)
		return
	}
	zapFields := convertToZapFields(fields)
	Logger.Warn(msg, zapFields...)
}

// Error logs a message at error level
func Error(msg string, fields ...LogField) {
	if len(fields) == 0 {
		Sugar.Error(msg)
		return
	}
	zapFields := convertToZapFields(fields)
	Logger.Error(msg, zapFields...)
}

// Fatal logs a message at fatal level and then exits
func Fatal(msg string, fields ...LogField) {
	if len(fields) == 0 {
		Sugar.Fatal(msg)
		return
	}
	zapFields := convertToZapFields(fields)
	Logger.Fatal(msg, zapFields...)
}

// WithFields creates a new logger with the given fields
func WithFields(fields ...LogField) *zap.Logger {
	return Logger.With(convertToZapFields(fields)...)
}

// WithError creates a new logger with an error field
func WithError(err error) *zap.Logger {
	return Logger.With(zap.Error(err))
}

// convertToZapFields converts LogField slice to zap.Field slice
func convertToZapFields(fields []LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
