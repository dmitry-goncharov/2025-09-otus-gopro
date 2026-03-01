package app

import (
	"context"
	"errors"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type ContextKey string

var UserIDKey ContextKey

var (
	ErrNoUserID      = errors.New("no user id")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrUnexpected    = errors.New("unexpected error")
)

type Application interface {
	CreateEvent(ctx context.Context, title string) (*Event, error)
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, evtID string) error
	GetDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, evt storage.Event) error
	UpdateEvent(ctx context.Context, evtID string, evt storage.Event) error
	DeleteEvent(ctx context.Context, evtID string) error
	GetDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsByRange(ctx context.Context, begin time.Time, end time.Time) ([]storage.Event, error)
	DeleteOutdatedEvents(ctx context.Context, date time.Time) error
}

type App struct {
	Logger  Logger
	Storage Storage
}

func NewApplication(logger Logger, storage Storage) Application {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, title string) (*Event, error) {
	a.Logger.Debug("create event, title:" + title)

	userID := ctx.Value(UserIDKey)
	if userID == "" {
		return nil, ErrNoUserID
	}

	evt := storage.Event{
		ID:     uuid.New().String(),
		Title:  title,
		Date:   time.Now(),
		UserID: userID.(string),
	}

	err := a.Storage.CreateEvent(ctx, evt)
	if err != nil {
		return nil, MapStorageErrToAppErr(err)
	}
	return MapStorageEventToAppEvent(&evt), nil
}

func (a *App) UpdateEvent(ctx context.Context, event *Event) error {
	a.Logger.Debug("update event, event:" + event.String())
	err := a.Storage.UpdateEvent(ctx, event.ID, *MapAppEventToStorageEvent(event))
	if err != nil {
		return MapStorageErrToAppErr(err)
	}
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, evtID string) error {
	a.Logger.Debug("delete event, ID:" + evtID)

	err := a.Storage.DeleteEvent(ctx, evtID)
	if err != nil {
		return MapStorageErrToAppErr(err)
	}
	return nil
}

func (a *App) GetDayEvents(ctx context.Context, date time.Time) ([]Event, error) {
	a.Logger.Debug("get day events, date:" + date.String())

	events, err := a.Storage.GetDayEvents(ctx, date)
	if err != nil {
		return nil, MapStorageErrToAppErr(err)
	}
	return MapStorageEventsToAppEvents(events), nil
}

func (a *App) GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error) {
	a.Logger.Debug("get week events, date:" + date.String())

	events, err := a.Storage.GetWeekEvents(ctx, date)
	if err != nil {
		return nil, MapStorageErrToAppErr(err)
	}
	return MapStorageEventsToAppEvents(events), nil
}

func (a *App) GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error) {
	a.Logger.Debug("get month events, date:" + date.String())

	events, err := a.Storage.GetMonthEvents(ctx, date)
	if err != nil {
		return nil, MapStorageErrToAppErr(err)
	}
	return MapStorageEventsToAppEvents(events), nil
}
