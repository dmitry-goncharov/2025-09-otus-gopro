package storage

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID     string    `db:"id"`
	Title  string    `db:"title"`
	Date   time.Time `db:"date"`
	UserID string    `db:"user_id"`
}

var (
	ErrAlreadyExists = errors.New("event with given id already exists")
	ErrNotFound      = errors.New("event with given id was not found")
)

func (e *Event) String() string {
	return fmt.Sprintf("Event{ID: %s, Title: %s, Date: %v, UserID: %s}", e.ID, e.Title, e.Date, e.UserID)
}
