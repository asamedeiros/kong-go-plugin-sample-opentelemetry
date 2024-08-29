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
	DebugWithContext(ctx context.Context, msg string)
	Fatal(msg string)
	FatalWithContext(ctx context.Context, msg string)
	Info(msg string)
	InfoWithContext(ctx context.Context, msg string)
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

func (c *log) Debug(msg string) {
	c.zapLogger.Debug(msg)
}

func (c *log) DebugWithContext(ctx context.Context, msg string) {
	c.zapLogger.Debug(msg, zap.Any("context", ctx))
}

func (c *log) Fatal(msg string) {
	c.zapLogger.Fatal(msg)
}

func (c *log) FatalWithContext(ctx context.Context, msg string) {
	c.zapLogger.Fatal(msg, zap.Any("context", ctx))
}

func (c *log) Info(msg string) {
	c.zapLogger.Info(msg)
}

func (c *log) InfoWithContext(ctx context.Context, msg string) {
	c.zapLogger.Info(msg, zap.Any("context", ctx))
}
