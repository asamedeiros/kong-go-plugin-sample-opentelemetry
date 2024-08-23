package main

import (
	"context"
	"fmt"
	"os"

	"github.com/asamedeiros/kong-go-sample-ddtrace/pkg/log"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	_log "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	tracer "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Kong/go-pdk/server"
)

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("sample-ddtrace-plugin"),
			semconv.ServiceVersion("0.1.0"),
			semconv.DeploymentEnvironment(os.Getenv("ENVIRONMENT")),
		))
}

func newTracerProvider(ctx context.Context, rsc *resource.Resource) (*trace.TracerProvider, tracer.Tracer) {
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		panic(err)
	}

	/* samplingRate := 1.0
	p := os.Getenv("KONG_TRACING_SAMPLING_RATE")
	if p != "" {
		samplingRate, _ = strconv.ParseFloat(p, 64)
	} */

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		/* trace.WithSampler(
			trace.ParentBased(trace.TraceIDRatioBased(samplingRate)),
		), */
		trace.WithResource(rsc),
	)

	otel.SetTracerProvider(tracerProvider)

	// Finally, set the tracer that can be used for this package.
	tracer := tracerProvider.Tracer("sample-ddtrace")

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider, tracer
}

func shutdownTrace(ctx context.Context, tracerProvider *trace.TracerProvider) {
	if err := tracerProvider.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func newLogProvider(ctx context.Context, rsc *resource.Resource) (*_log.LoggerProvider, log.Log) {

	// Use a working LoggerProvider implementation instead e.g. use go.opentelemetry.io/otel/sdk/log.
	//provider := noop.NewLoggerProvider()

	// Create a logger provider.
	// You can pass this instance directly when creating bridges.
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		panic(err)
	}
	processor := _log.NewBatchProcessor(exporter)
	loggerProvider := _log.NewLoggerProvider(
		_log.WithResource(rsc),
		_log.WithProcessor(processor),
	)

	// Register as global logger provider so that it can be accessed global.LoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(loggerProvider)

	// Initialize a zap zaplogger with the otelzap bridge core.
	// This method actually doesn't log anything on your STDOUT, as everything
	// is shipped to a configured otel endpoint.
	zaplogger := zap.New(otelzap.NewCore("sample-ddtrace", otelzap.WithLoggerProvider(loggerProvider)))
	zap.ReplaceGlobals(zaplogger)

	//zapLog, _ := zap.NewProduction(zap.AddCallerSkip(1))
	//olog := otelzap.New(zapLog)
	sugar := zaplogger.Sugar()

	/* _log := _otelzap.New(zap.NewExample())
	_sugar := _log.Sugar()

	_sugar.Ctx(ctx).Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"attempt", 3,
		"backoff", time.Second,
	) */

	return loggerProvider, log.New(sugar)
}

func shutdownLogProvider(ctx context.Context, loggerProvider *_log.LoggerProvider, l log.Log) {

	l.Sync()
	/* if err := l.Sync(); err != nil {
		l.Warn("sync remaining logs failed: ", err.Error())
	} */

	if err := loggerProvider.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func main() {

	// Create resource.
	rsc, err := newResource()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tracerProvider, tracer := newTracerProvider(ctx, rsc)
	defer shutdownTrace(ctx, tracerProvider)

	logProvider, l := newLogProvider(ctx, rsc)
	defer shutdownLogProvider(ctx, logProvider, l)

	ctor := func() interface{} { return plugin.NewPlugin(l, tracer) }
	err = server.StartServer(ctor, plugin.Version, plugin.Priority)
	if err != nil {
		l.Error(fmt.Errorf("embedded plugin server start error: %w", err))
	}
}
