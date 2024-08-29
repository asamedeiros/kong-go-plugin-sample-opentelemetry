package log

import (
	"context"

	"go.uber.org/zap"
)

type Log interface {
	/* Level() zapcore.Level
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) Log */
	With(name, value string) Log
	Sync() error
	Error(msg string)
	ErrorWithContext(ctx context.Context, msg string)
	Warn(msg string)
	WarnWithContext(ctx context.Context, msg string)
}

func New(zapLogger *zap.Logger) Log {
	return &log{
		zapLogger: zapLogger,
	}
}

type log struct {
	zapLogger *zap.Logger
}

func (c *log) With(name, value string) Log {
	return New(c.zapLogger.With(zap.String(name, value)))
}

func (c *log) Sync() error {
	return c.zapLogger.Sync()
}

func (c *log) Error(msg string) {
	c.zapLogger.Error(msg)
}

func (c *log) ErrorWithContext(ctx context.Context, msg string) {
	c.zapLogger.Error(msg, zap.Any("context", ctx))
}

func (c *log) Warn(msg string) {
	c.zapLogger.Warn(msg)
}

func (c *log) WarnWithContext(ctx context.Context, msg string) {
	c.zapLogger.Warn(msg, zap.Any("context", ctx))
}

/* func (l *log) Level() zapcore.Level {
	return l.z.Level()
}

func (l *log) With(args ...interface{}) Log {
	return New(l.z.With(args...))
}

func (l *log) Debug(args ...interface{}) {
	l.z.Debug(args...)
}

func (l *log) Debugf(template string, args ...interface{}) {
	l.z.Debugf(template, args...)
}

func (l *log) Info(args ...interface{}) {
	l.z.Info(args...)
}

func (l *log) Infof(template string, args ...interface{}) {
	l.z.Infof(template, args...)
}

func (l *log) Warn(args ...interface{}) {
	l.z.Warn(args...)
}

func (l *log) Warnf(template string, args ...interface{}) {
	l.z.Warnf(template, args...)
}

func (l *log) Error(args ...interface{}) {
	l.z.Error(args...)
}

func (l *log) Errorf(template string, args ...interface{}) {
	l.z.Errorf(template, args...)
}

func (l *log) Fatal(args ...interface{}) {
	l.z.Fatal(args...)
} */
