package internalhttp

import (
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

const (
	ID     = "id"
	Period = "p"
	Day    = "d"
	Week   = "w"
	Month  = "m"
	Year   = "y"
)

type EventTitleDto struct {
	Title string `json:"title"`
}

type EventDto struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserID string    `json:"userId"`
}

func MapAppEventToEventDto(appEvent *app.Event) *EventDto {
	return &EventDto{
		ID:     appEvent.ID,
		Title:  appEvent.Title,
		Date:   appEvent.Date,
		UserID: appEvent.UserID,
	}
}

func MapEventDtoToAppEvent(eventDto *EventDto) *app.Event {
	return &app.Event{
		ID:     eventDto.ID,
		Title:  eventDto.Title,
		Date:   eventDto.Date,
		UserID: eventDto.UserID,
	}
}

func MapAppEventsToEventDtos(appEvents []app.Event) []EventDto {
	eventDtos := make([]EventDto, len(appEvents))
	for i, appEvent := range appEvents {
		eventDtos[i] = *MapAppEventToEventDto(&appEvent)
	}
	return eventDtos
}
