package apm

import (
	"context"
	"fmt"
	"os"
	"strings"

	internalLog "github.com/asamedeiros/kong-go-sample-ddtrace/internal/log"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
)

var shutdowns []func()

func newResource(name, version, environment string) (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = uuid.New().String()
	}
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(name),
			semconv.ServiceVersion(version),
			semconv.DeploymentEnvironment(environment),
			semconv.ServiceInstanceIDKey.String(hostname),
		))
}

func configLog(ctx context.Context, rsc *resource.Resource) (internalLog.Log, func(), error) {

	// Create a logger provider.
	// You can pass this instance directly when creating bridges.
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		// Criando um log emergencial, para conseguir reportar errors, já que não chegará a instanciar o oficial mais abaixo
		l, _ := zap.NewProduction(zap.AddCallerSkip(1))
		return internalLog.New(l), nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(rsc),
		log.WithProcessor(processor),
	)

	// Register as global logger provider so that it can be accessed global.LoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(loggerProvider)

	// Initialize a zap zaplogger with the otelzap bridge core.
	// This method actually doesn't log anything on your STDOUT, as everything
	// is shipped to a configured otel endpoint.
	zaplogger := zap.New(otelzap.NewCore(plugin.Name, otelzap.WithLoggerProvider(loggerProvider)))

	/* // Wrap zap logger to extend Zap with API that accepts a context.Context.
	zaploggerWrap := _otelzap.New(zaplogger)
	sugarWrap := zaploggerWrap.Sugar()

	return loggerProvider, log.New(sugarWrap), nil */

	log := internalLog.New(zaplogger)

	shutdown := func() {
		if err := log.Sync(); err != nil && !strings.Contains(err.Error(), "sync /dev/stdout: invalid argument") {
			log.Warn("error syncing logger: " + err.Error())
		}

		if loggerProvider != nil {
			if err := loggerProvider.Shutdown(ctx); err != nil {
				log.Warn("error shutting down logger provider: " + err.Error())
			}
		}
	}

	return log, shutdown, nil
}

func configTracerProvider(ctx context.Context, l internalLog.Log, rsc *resource.Resource) (func(), error) {
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
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
	//tracer := otel.GetTracerProvider().Tracer("sample-ddtrace")

	//otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	shutdown := func() {
		if tracerProvider != nil {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				l.Warn("error shutting down tracer provider: " + err.Error())
			}
		}
	}

	return shutdown, nil
}

func ConfigOpenTelemetry(name, version, environment string) internalLog.Log {

	rsc, err := newResource(name, version, environment)
	if err != nil {
		fmt.Printf("error in plugin - failed to create resource for log provider: %s", err)
	}

	ctx := context.Background()

	log, shutdownLog, err := configLog(ctx, rsc)
	if err != nil {
		log.Error(fmt.Sprintf("error in plugin - failed to create log or log provider: %s", err))
	}
	shutdowns = append(shutdowns, shutdownLog)

	shutdownTracerProvider, err := configTracerProvider(ctx, log, rsc)
	if err != nil {
		log.Error(fmt.Sprintf("error in plugin - failed to create trace provider: %s", err))
	}
	shutdowns = append(shutdowns, shutdownTracerProvider)

	return log
}

func StopOpenTelemetry() {
	for _, s := range shutdowns {
		s()
	}
}
