package plugin

import (
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/asamedeiros/kong-go-sample-ddtrace/pkg/log"
)

const (
	Version  = "0.2"
	Priority = 1400
)

type Config interface {
	Access(kong *pdk.PDK)
}

type pluginConfig struct {
	log log.Log
}

// NewPlugin returns a new plugin configuration.
func NewPlugin(log log.Log) Config {

	return &pluginConfig{
		log: log,
	}
}

// Access is executed for every request from a client
// and, before it is being proxied to the upstream service.
func (c *pluginConfig) Access(kong *pdk.PDK) {

	//c.log.Error(fmt.Sprintf("error_2 - %s", "opa"))

	//c.log.Info(fmt.Sprintf("info_2 - %s", "opa"))

	str, _ := kong.Log.Serialize()

	kong.Log.Err("error_kong_2: ", str)

	kong.Log.Info("info_kong_2")

	//kong.Log.Err("error_kong_3, a: b, f: d")

	//c.accessError(kong, 401)

	/* h, _ := kong.Request.GetHeaders(-1)
	rHeader := make(map[string]string)
	for k := range h {
		rHeader[strings.ToLower(k)] = h[k][0]
	}

	rPath, _ := kong.Request.GetPath()
	rMethod, _ := kong.Request.GetMethod()
	rHost, _ := kong.Request.GetHost()
	rRawQuery, _ := kong.Request.GetRawQuery()
	rRemoteAddr, _ := kong.Client.GetIp()
	rBody, _ := kong.Request.GetRawBody()

	req := &entities.PermissionRequest{
		Header:  rHeader,
		Method:  rMethod,
		RawBody: rBody,
		URL: &url.URL{
			Host:     rHost,
			Path:     rPath,
			RawQuery: rRawQuery,
		},
		RemoteAddr: rRemoteAddr,
	}

	reqLog := c.log.With("plugin", "sample-ddtrace").
		With("x-request-id", req.GetHeader("x-request-id")).
		With("method", req.Method).
		With("path", req.URL.Path).
		With("host", req.GetHeader("host")).
		With("user-agent", req.GetHeader("user-agent")).
		With("remote-addr", req.RemoteAddr).
		With("cf-ray", req.GetHeader("cf-ray")).
		With("aws-xray", req.GetHeader("x-amzn-trace-id"))

	reqLog.Error(fmt.Sprintf("error_3 - %s", "opa")) */

	//c.accessError(kong, rsl.Status)
}

func (c *pluginConfig) accessError(kong *pdk.PDK, code int) {
	headers := make(map[string][]string)
	if code == http.StatusUnauthorized {
		kong.Response.AddHeader("X-sample-ddtrace", "Unauthorized")
		kong.Response.Exit(code, []byte("Unauthorized"), headers)
	} else {
		kong.Response.AddHeader("X-sample-ddtrace", "Forbidden")
		kong.Response.Exit(code, []byte("Forbidden"), headers)
	}
}
