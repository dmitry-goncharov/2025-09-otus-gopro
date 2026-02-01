package internalhttp

import (
	"net/http"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Handler struct {
	log app.Logger
	app app.Application
}

func NewHandler(log app.Logger, app app.Application) *Handler {
	return &Handler{
		log: log,
		app: app,
	}
}

func (h *Handler) Hello(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(`hello-world`))
	if err != nil {
		h.log.Error("error handling hello")
	}
}
