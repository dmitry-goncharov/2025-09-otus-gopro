package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	ID     string
	Title  string
	Date   time.Time
	UserID string
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{ID: %s, Title: %s, Date: %v, UserID: %s}", e.ID, e.Title, e.Date, e.UserID)
}

func MapAppEventToStorageEvent(appEvent *Event) *storage.Event {
	return &storage.Event{
		ID:     appEvent.ID,
		Title:  appEvent.Title,
		Date:   appEvent.Date,
		UserID: appEvent.UserID,
	}
}

func MapStorageEventToAppEvent(storageEvent *storage.Event) *Event {
	return &Event{
		ID:     storageEvent.ID,
		Title:  storageEvent.Title,
		Date:   storageEvent.Date,
		UserID: storageEvent.UserID,
	}
}

func MapStorageEventsToAppEvents(storageEvents []storage.Event) []Event {
	appEvents := make([]Event, len(storageEvents))
	for i, storageEvent := range storageEvents {
		appEvents[i] = *MapStorageEventToAppEvent(&storageEvent)
	}
	return appEvents
}

func MapStorageErrToAppErr(storageError error) error {
	switch {
	case errors.Is(storageError, storage.ErrAlreadyExists):
		return errors.Join(ErrAlreadyExists, storageError)
	case errors.Is(storageError, storage.ErrNotFound):
		return errors.Join(ErrNotFound, storageError)
	default:
		return errors.Join(ErrUnexpected, storageError)
	}
}
