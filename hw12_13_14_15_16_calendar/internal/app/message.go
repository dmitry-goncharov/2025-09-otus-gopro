package app

import (
	"fmt"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
)

type Message struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserID string    `json:"userId"`
}

func (e *Message) String() string {
	return fmt.Sprintf("Message{ID: %s, Title: %s, Date: %v, UserID: %s}", e.ID, e.Title, e.Date, e.UserID)
}

func MapStorageEventToMessage(event *storage.Event) *Message {
	return &Message{
		ID:     event.ID,
		Title:  event.Title,
		Date:   event.Date,
		UserID: event.UserID,
	}
}

func MapStorageEventsToMessages(events []storage.Event) []Message {
	messages := make([]Message, len(events))
	for i, event := range events {
		messages[i] = *MapStorageEventToMessage(&event)
	}
	return messages
}
