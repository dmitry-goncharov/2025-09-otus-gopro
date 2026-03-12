package scanner

import (
	"context"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Scanner interface {
	Scan(ctx context.Context) []app.Message
}

type StorageScanner struct {
	log      app.Logger
	storage  app.Storage
	interval time.Duration
}

func NewStorageScanner(log app.Logger, storage app.Storage, interval time.Duration) Scanner {
	return &StorageScanner{
		log:      log,
		storage:  storage,
		interval: interval,
	}
}

func (s *StorageScanner) Scan(ctx context.Context) []app.Message {
	s.log.Debug("scan storage for interval " + s.interval.String())
	events, err := s.storage.GetEventsByRange(ctx, time.Now(), time.Now().Add(s.interval))
	if err != nil {
		s.log.Error("error scanning storage for interval " + s.interval.String())
	}
	return app.MapStorageEventsToMessages(events)
}
