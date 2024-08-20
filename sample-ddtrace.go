package main

import (
	"log"

	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"

	"github.com/Kong/go-pdk/server"
)

func main() {

	ctor := func() interface{} { return plugin.NewPlugin() }
	err := server.StartServer(ctor, plugin.Version, plugin.Priority)
	if err != nil {
		log.Printf("embedded plugin server start error: %w", err)
	}
}
