package log

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	z *zap.SugaredLogger
}

// NewZapLogger creates a custom zap.SugaredLogger implementation of the Logger interface
func NewZapLogger(config zap.Config, opts ...zap.Option) (Logger, error) {
	logger, err := config.Build(opts...)
	return &zapLogger{
		z: logger.Sugar(),
	}, err
}

// NewNoOpZapLogger creates a NoOp Logger for testing purposes
func NewNoOpZapLogger() Logger {
	logger := zap.NewNop()
	return &zapLogger{
		z: logger.Sugar(),
	}
}

// Error logs a message at the error level
func (l *zapLogger) Error(args ...interface{}) {
	l.z.Error(args...)
}

// Info logs a message at the info level
func (l *zapLogger) Info(args ...interface{}) {
	l.z.Info(args...)
}

// Debug logs a message at the debug level
func (l *zapLogger) Debug(args ...interface{}) {
	l.z.Debug(args...)
}

// Errorf logs a message at the error level
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.z.Errorf(template, args...)
}

// Infof logs a message at the info level
func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.z.Infof(template, args...)
}

// Debugf logs a message at the debug level
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.z.Debugf(template, args...)
}

// With returns a new a new zap.SugaredLogger Logger with the additional args as key-value context
func (l *zapLogger) With(args ...interface{}) Logger {
	logger := l.z.With(args...)
	return &zapLogger{
		z: logger,
	}
}
