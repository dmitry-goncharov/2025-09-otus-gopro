package app

import (
	"context"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

const (
	UserID = "UserId"
)

type Application interface {
	CreateEvent(ctx context.Context, title string) error
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

func (a *App) CreateEvent(ctx context.Context, title string) error {
	a.Logger.Debug("Create event, title:" + title)
	evt := storage.Event{
		ID:     uuid.New().String(),
		Title:  title,
		Date:   time.Now(),
		UserID: ctx.Value(UserID).(string),
	}
	return a.Storage.CreateEvent(ctx, evt)
}
