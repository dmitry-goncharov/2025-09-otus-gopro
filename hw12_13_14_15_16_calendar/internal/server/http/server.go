package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
)

const (
	UserID = "X-User-Id"
)

type Server struct {
	log app.Logger
	srv *http.Server
}

func NewServer(conf config.ServerConf, log app.Logger, app app.Application) *Server {
	handler := NewHandler(log, app)

	router := http.NewServeMux()
	router.HandleFunc("GET /hello", handler.Hello)
	router.HandleFunc("POST /events", handler.AddEvent)
	router.HandleFunc("POST /events/{id}", handler.UpdateEvent)
	router.HandleFunc("DELETE /events/{id}", handler.DeleteEvent)
	router.HandleFunc("GET /events", handler.GetEvents)

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
	s.log.Info("http server starting on addr: " + s.srv.Addr)
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("http server stopping")
	return s.srv.Shutdown(ctx)
}
