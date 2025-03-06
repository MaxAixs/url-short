package server

import (
	"context"
	"net/http"
	"url-shortener/app/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.HTTPServer, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
