package apm

import (
	"context"

	"github.com/Kong/go-pdk"
	"github.com/asamedeiros/kong-go-sample-ddtrace/internal/log"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Wrapper interface {
	Access(kong *pdk.PDK)
}

type wrapperConfig struct {
	log    log.Log
	plugin plugin.Plugin
}

func NewPluginWrapper(log log.Log, plugin plugin.Plugin) Wrapper {

	return &wrapperConfig{
		log:    log,
		plugin: plugin,
	}
}

func (c *wrapperConfig) wrapper(ctx context.Context, kong *pdk.PDK) (context.Context, log.Log) {

	logWithAtts := c.log

	requestid, _ := kong.Request.GetHeader("x-request-id")
	if requestid != "" {
		logWithAtts = logWithAtts.With("request.x_request_id", requestid)
	}

	// primeiramente tenta recuperar o traceparent do contexto compartilhado (setado pelo kong-plugin-tracing-customizations)
	traceparent, _ := kong.Ctx.GetSharedString("traceparent")
	if traceparent == "" {
		logWithAtts.Infof("getting traceparent from header: %s", traceparent)
		traceparent, _ = kong.Request.GetHeader("traceparent")
	} else {
		logWithAtts.Infof("getting traceparent from shared: %s", traceparent)
	}
	if traceparent != "" {
		// prepare carrier to set traceparent into context
		carrier := propagation.MapCarrier{}
		carrier.Set("traceparent", traceparent)
		// reads tracecontext from the carrier into a returned Context.
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	}

	return ctx, logWithAtts
}

func (c *wrapperConfig) Access(kong *pdk.PDK) {
	ctx, logWithAtts := c.wrapper(context.Background(), kong)
	c.plugin.Access(ctx, logWithAtts, kong)
}
