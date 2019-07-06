package server

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/sirupsen/logrus"
	cloki "github.com/ssfilatov/clickhouse-loki-adapter/core"
	"github.com/ssfilatov/clickhouse-loki-adapter/core/clickhouse"
	"github.com/ssfilatov/clickhouse-loki-adapter/core/querier"
	"net"
	"net/http"
)

type Server struct {
	cfg          *cloki.Config
	httpListener net.Listener

	HTTP       *mux.Router
	HTTPServer *http.Server
}

func New(cfg *cloki.Config) (*Server, error) {
	// Setup listeners first, so we can fail early if the port is in use.
	httpListener, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%d", cfg.Server.HTTPListenHost, cfg.Server.HTTPListenPort))
	if err != nil {
		return nil, err
	}

	ch, err := clickhouse.NewClickhouse(cfg)
	if err != nil {
		return nil, err
	}
	q := querier.New(cfg, ch)

	router := mux.NewRouter()
	router.Handle("/api/prom/query", http.HandlerFunc(q.QueryHandler))
	//add these queries later
	router.Handle("/api/prom/label", http.HandlerFunc(q.LabelHandler))
	router.Handle("/api/prom/label/{name}/values", http.HandlerFunc(q.LabelHandler))

	httpServer := &http.Server{
		ReadTimeout:  cfg.Server.HTTPServerReadTimeout,
		WriteTimeout: cfg.Server.HTTPServerWriteTimeout,
		IdleTimeout:  cfg.Server.HTTPServerIdleTimeout,
		Handler:      router,
	}

	return &Server{
		cfg:          cfg,
		httpListener: httpListener,

		HTTP:       router,
		HTTPServer: httpServer,
	}, nil
}

func (s *Server) Run() error {

	err := s.HTTPServer.Serve(s.httpListener)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
