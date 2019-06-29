package chloki

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gitlab.corp.mail.ru/infra/dpp-server/core/clickhouse"
	"time"
)

var (
	config             *Config
)

func defaultHandler(ctx *fasthttp.RequestCtx) {
	ctx.Error("Not implemented", fasthttp.StatusNotImplemented)
}

func Run(configFileName string) error {

	config = ParseConfig(configFileName)

	router := fasthttprouter.New()
	router.GET("/api/prom/query", queryHandler)
	router.GET("/api/prom/label", getLabelsHandler)
	router.GET("/api/prom/label/:name/values", getLabelValuesHandler)

	srv := &fasthttp.Server{
		MaxRequestBodySize: 16 << 20,
		Handler:            router.Handler,
	}

	clickhouse.InsertProxyClient = &fasthttp.HostClient{
		Addr:                config.Clickhouse.Server,
		MaxConns:            config.Clickhouse.MaxConnections,
		MaxIdleConnDuration: time.Second,
	}

	if config.UseTls {
		return srv.ListenAndServeTLS(config.ListenAddr, config.CertFile, config.KeyFile)
	} else {
		return srv.ListenAndServe(config.ListenAddr)
	}

}
