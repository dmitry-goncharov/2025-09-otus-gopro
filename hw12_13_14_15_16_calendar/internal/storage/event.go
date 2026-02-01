package storage

import (
	"errors"
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
