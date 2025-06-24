package logx

import (
	"go.uber.org/zap/zapcore"
)

// A Logger represents a logger.
type Logger interface {
	// Debug logs a message at info level.
	Debug(...any)
	// Debugf logs a message at info level.
	Debugf(string, ...any)
	// Debugv logs a message at info level.
	Debugv(any)
	// Debugw logs a message at info level.
	Debugw(string, ...zapcore.Field)
	// Error logs a message at error level.
	Error(...any)
	// Errorf logs a message at error level.
	Errorf(string, ...any)
	// Errorv logs a message at error level.
	Errorv(any)
	// Errorw logs a message at error level.
	Errorw(string, ...zapcore.Field)
	// Warn logs a message at warn level.
	Warn(...any)
	// Warnf logs a message at warn level.
	Warnf(string, ...any)
	// Warnv logs a message at warn level.
	Warnv(any)
	// Warnw logs a message at warn level.
	Warnw(string, ...zapcore.Field)
	// Info logs a message at info level.
	Info(...any)
	// Infof logs a message at info level.
	Infof(string, ...any)
	// Infov logs a message at info level.
	Infov(any)
	// Infow logs a message at info level.
	Infow(string, ...zapcore.Field)
}
