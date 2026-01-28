package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	log app.Logger
	srv *http.Server
}

func NewServer(conf config.ServerConf, log app.Logger, app app.Application) *Server {
	handler := NewHandler(log, app)

	router := http.NewServeMux()
	router.HandleFunc("GET /hello", handler.Hello)

	server := &http.Server{
		Addr:              net.JoinHostPort(conf.Host, conf.Port),
		Handler:           loggingMiddleware(log, router),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &Server{
		log: log,
		srv: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("Server starting on addr: " + s.srv.Addr)
	err := s.srv.ListenAndServe()
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Server stopping")
	return s.srv.Shutdown(ctx)
}
