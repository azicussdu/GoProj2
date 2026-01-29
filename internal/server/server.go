package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	myServer *http.Server
}

func New(router http.Handler, port string) *Server {
	return &Server{
		myServer: &http.Server{
			Addr:    ":" + port, // "8080"
			Handler: router,
			// TODO configs
		},
	}
}

func (s *Server) Run() error {
	slog.Info("server is started", "ADDR", s.myServer.Addr)
	return s.myServer.ListenAndServe()
}
