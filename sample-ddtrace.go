package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/asamedeiros/kong-go-sample-ddtrace/pkg/log"
	"github.com/asamedeiros/kong-go-sample-ddtrace/plugin"

	"github.com/Kong/go-pdk/server"
)

func setupLog() log.Log {
	zapLog, _ := zap.NewProduction(zap.AddCallerSkip(1))
	sugar := zapLog.Sugar()

	return log.New(sugar)
}

func syncLog(l log.Log) {
	l.Sync()
	/* if err := l.Sync(); err != nil {
		l.Warn("sync remaining logs failed: ", err.Error())
	} */
}

func main() {
	l := setupLog()
	defer syncLog(l)

	//l.Error("erro_1")
	//l.Info("info_1")

	ctor := func() interface{} { return plugin.NewPlugin() }
	err := server.StartServer(ctor, plugin.Version, plugin.Priority)
	if err != nil {
		l.Error(fmt.Errorf("embedded plugin server start error: %w", err))
	}
}
