package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) AddEvent(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("add event handler")

	userID := r.Header.Get(UserID)
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dto := EventTitleDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.log.Error("error decoding event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), app.UserIDKey, userID), time.Second*5)
	defer cancel()

	event, err := h.app.CreateEvent(ctx, dto.Title)
	if err != nil {
		h.logAppError(err, "create event")
		w.WriteHeader(mapAppErrorToHTTPStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(MapAppEventToEventDto(event))
	if err != nil {
		h.log.Error("error ecoding event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("update event handler")

	id := r.PathValue(ID)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dto := EventDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.log.Error("error decoding event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if id != dto.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = h.app.UpdateEvent(ctx, MapEventDtoToAppEvent(&dto))
	if err != nil {
		h.logAppError(err, "update event")
		w.WriteHeader(mapAppErrorToHTTPStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		h.log.Error("error ecoding event: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("delete event handler")

	id := r.PathValue(ID)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := h.app.DeleteEvent(ctx, id)
	if err != nil {
		h.logAppError(err, "delete event")
		w.WriteHeader(mapAppErrorToHTTPStatus(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("get events handler")

	period := r.URL.Query().Get(Period)
	if period == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(r.URL.Query().Get(Year))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(r.URL.Query().Get(Month))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	day, err := strconv.Atoi(r.URL.Query().Get(Day))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	var events []app.Event
	switch period {
	case Day:
		events, err = h.app.GetDayEvents(ctx, date)
		if err != nil {
			h.logAppError(err, "get day events")
			w.WriteHeader(mapAppErrorToHTTPStatus(err))
			return
		}
	case Week:
		events, err = h.app.GetWeekEvents(ctx, date)
		if err != nil {
			h.logAppError(err, "get week events")
			w.WriteHeader(mapAppErrorToHTTPStatus(err))
			return
		}
	case Month:
		events, err = h.app.GetMonthEvents(ctx, date)
		if err != nil {
			h.logAppError(err, "get month events")
			w.WriteHeader(mapAppErrorToHTTPStatus(err))
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(MapAppEventsToEventDtos(events))
	if err != nil {
		h.log.Error("error ecoding events: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) logAppError(appError error, action string) {
	switch {
	case errors.Is(appError, app.ErrNoUserID):
		h.log.Debug("no user id at " + action + ", error:" + appError.Error())
	case errors.Is(appError, app.ErrAlreadyExists):
		h.log.Debug("event already exists at " + action + ", error:" + appError.Error())
	case errors.Is(appError, app.ErrNotFound):
		h.log.Debug("event not found at " + action + ", error:" + appError.Error())
	case errors.Is(appError, app.ErrUnexpected):
		h.log.Error("unexpected error at " + action + ", error:" + appError.Error())
	default:
		h.log.Error("undefined error at " + action + ", error:" + appError.Error())
	}
}

func mapAppErrorToHTTPStatus(appError error) int {
	switch {
	case errors.Is(appError, app.ErrNoUserID):
		return http.StatusBadRequest
	case errors.Is(appError, app.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(appError, app.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(appError, app.ErrUnexpected):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
