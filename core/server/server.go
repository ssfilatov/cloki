package server

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/sirupsen/logrus"
	cloki "github.com/ssfilatov/cloki/core"
	"net"
	"net/http"
)

type Server struct {
	httpListener net.Listener

	HTTP       *mux.Router
	HTTPServer *http.Server
}

func New(cloki *cloki.CLoki) (*Server, error) {
	// Setup listeners first, so we can fail early if the port is in use.
	httpListener, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%d", cloki.Cfg.Server.HTTPListenHost, cloki.Cfg.Server.HTTPListenPort))
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	router.Handle("/api/prom/query", http.HandlerFunc(cloki.Query.QueryHandler))
	//add these queries later
	router.Handle("/api/prom/label", http.HandlerFunc(cloki.Query.LabelHandler))
	router.Handle("/api/prom/label/{name}/values", http.HandlerFunc(cloki.Query.LabelHandler))

	httpServer := &http.Server{
		ReadTimeout:  cloki.Cfg.Server.HTTPServerReadTimeout,
		WriteTimeout: cloki.Cfg.Server.HTTPServerWriteTimeout,
		IdleTimeout:  cloki.Cfg.Server.HTTPServerIdleTimeout,
		Handler:      router,
	}

	return &Server{
		httpListener: httpListener,
		HTTP:         router,
		HTTPServer:   httpServer,
	}, nil
}

func (s *Server) Run() error {

	err := s.HTTPServer.Serve(s.httpListener)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
