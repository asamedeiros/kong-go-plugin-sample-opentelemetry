package plugin

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/Kong/go-pdk"
	"github.com/asamedeiros/kong-go-sample-ddtrace/internal/log"
)

const (
	Name     = "kong-go-plugin-sample-ddtrace"
	Version  = "0.2"
	Priority = 1000
)

type Plugin interface {
	Access(ctx context.Context, log log.Log, kong *pdk.PDK)
}

type pluginConfig struct {
}

// NewPlugin returns a new plugin configuration.
func NewPlugin() Plugin {

	return &pluginConfig{}
}

// Access is executed for every request from a client
// and, before it is being proxied to the upstream service.
func (c *pluginConfig) Access(ctx context.Context, log log.Log, kong *pdk.PDK) {

	log = log.With("a", "b").With("c", "d").With("e", "f")

	tracer := otel.GetTracerProvider().Tracer("sample-ddtrace")
	_, span2 := tracer.Start(ctx, "access")
	defer span2.End()

	log.Error("setting sample")

	log.ErrorWithContext(ctx, "setting sample with context")
}
