package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	httpSrv *http.Server
}

func New(router http.Handler, port string) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:    ":" + port, // "8080"
			Handler: router,
			// TODO configs
		},
	}
}

func (s *Server) Run() error {
	slog.Info("server is started", "ADDR", s.httpSrv.Addr)
	return s.httpSrv.ListenAndServe()
}
