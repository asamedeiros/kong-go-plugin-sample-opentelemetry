package main

import (
	"fmt"
	"os"

	"github.com/asamedeiros/kong-go-sample-ddtrace/internal/apm"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"

	"github.com/Kong/go-pdk/server"
)

func main() {

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "staging"
	}

	log := apm.ConfigOpenTelemetry(plugin.Name, plugin.Version, environment)
	defer apm.StopOpenTelemetry()

	ctor := func() interface{} { return apm.NewPluginWrapper(log, plugin.NewPlugin()) }

	err := server.StartServer(ctor, plugin.Version, plugin.Priority)
	if err != nil {
		log.Error(fmt.Sprintf("embedded plugin server start error: %s", err))
	}
}
