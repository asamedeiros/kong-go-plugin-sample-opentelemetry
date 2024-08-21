package main

import (
	"context"
	"fmt"

	"github.com/asamedeiros/kong-go-sample-ddtrace/pkg/log"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	_log "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"

	"github.com/Kong/go-pdk/server"
)

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("sample-ddtrace-plugin"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func newLoggerProvider(ctx context.Context, res *resource.Resource) (*_log.LoggerProvider, error) {
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}
	processor := _log.NewBatchProcessor(exporter)
	provider := _log.NewLoggerProvider(
		_log.WithResource(res),
		_log.WithProcessor(processor),
	)
	return provider, nil
}

func setupLog(ctx context.Context) (*_log.LoggerProvider, log.Log) {

	// Use a working LoggerProvider implementation instead e.g. use go.opentelemetry.io/otel/sdk/log.
	//provider := noop.NewLoggerProvider()

	// Create resource.
	res, err := newResource()
	if err != nil {
		panic(err)
	}

	// Create a logger provider.
	// You can pass this instance directly when creating bridges.
	loggerProvider, err := newLoggerProvider(ctx, res)
	if err != nil {
		panic(err)
	}

	// Register as global logger provider so that it can be accessed global.LoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(loggerProvider)

	// Initialize a zap zaplogger with the otelzap bridge core.
	// This method actually doesn't log anything on your STDOUT, as everything
	// is shipped to a configured otel endpoint.
	zaplogger := zap.New(otelzap.NewCore("sample-ddtrace", otelzap.WithLoggerProvider(loggerProvider)))

	//zapLog, _ := zap.NewProduction(zap.AddCallerSkip(1))
	//olog := otelzap.New(zapLog)
	sugar := zaplogger.Sugar()

	return loggerProvider, log.New(sugar)
}

func syncLog(ctx context.Context, loggerProvider *_log.LoggerProvider, l log.Log) {

	l.Sync()
	/* if err := l.Sync(); err != nil {
		l.Warn("sync remaining logs failed: ", err.Error())
	} */

	if err := loggerProvider.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func main() {

	ctx := context.Background()

	logProvider, l := setupLog(ctx)
	defer syncLog(ctx, logProvider, l)

	ctor := func() interface{} { return plugin.NewPlugin(l) }
	err := server.StartServer(ctor, plugin.Version, plugin.Priority)
	if err != nil {
		l.Error(fmt.Errorf("embedded plugin server start error: %w", err))
	}
}
