package log

import (
	"context"
	"fmt"

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
	WithTracing(ctx context.Context) Log
	Sync() error
	Error(msg string)
	Errorf(msg string, args ...interface{})
	Warn(msg string)
	Fatal(msg string)
	Info(msg string)
	Infof(msg string, args ...interface{})
}

func New(zapLogger *zap.Logger) Log {
	return &log{
		zapLogger: zapLogger,
	}
}

type log struct {
	zapLogger *zap.Logger
}

func (c *log) with(field zap.Field) Log {
	return New(c.zapLogger.With(field))
}

func (c *log) With(name, value string) Log {
	return c.with(zap.String(name, value))
}

func (c *log) WithTracing(ctx context.Context) Log {
	return c.with(zap.Any("context", ctx))
}

func (c *log) Sync() error {
	return c.zapLogger.Sync()
}

func (c *log) Error(msg string) {
	c.zapLogger.Error(msg)
}

func (c *log) Errorf(msg string, args ...interface{}) {
	c.zapLogger.Error(fmt.Sprintf(msg, args...))
}

func (c *log) Warn(msg string) {
	c.zapLogger.Warn(msg)
}

func (c *log) Debug(msg string) {
	c.zapLogger.Debug(msg)
}

func (c *log) Fatal(msg string) {
	c.zapLogger.Fatal(msg)
}

func (c *log) Info(msg string) {
	c.zapLogger.Info(msg)
}

func (c *log) Infof(msg string, args ...interface{}) {
	c.zapLogger.Info(fmt.Sprintf(msg, args...))
}
