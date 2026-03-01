package memorystorage

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu   sync.RWMutex
	evts map[string]storage.Event
}

func New() app.Storage {
	return &Storage{
		mu:   sync.RWMutex{},
		evts: make(map[string]storage.Event),
	}
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, evt storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.evts[evt.ID]; ok {
		return storage.ErrAlreadyExists
	}

	s.evts[evt.ID] = evt

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, evtID string, evt storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.evts[evtID]; !ok {
		return storage.ErrNotFound
	}

	s.evts[evt.ID] = evt

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, evtID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.evts[evtID]; !ok {
		return storage.ErrNotFound
	}

	delete(s.evts, evtID)

	return nil
}

func (s *Storage) GetDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 1)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 7)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 1, 0)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetEventsByRange(_ context.Context, begin time.Time, end time.Time) ([]storage.Event, error) {
	res := make([]storage.Event, 0)
	for _, evt := range s.evts {
		slog.Debug("get events by range", slog.Any("evtDate", evt.Date), slog.Any("begin", begin), slog.Any("end", end))
		ok := (evt.Date.Equal(begin) || evt.Date.After(begin)) && evt.Date.Before(end)
		if ok {
			res = append(res, evt)
		}
	}
	return res, nil
}

func (s *Storage) DeleteOutdatedEvents(_ context.Context, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, e := range s.evts {
		if e.Date.Before(date) {
			delete(s.evts, key)
		}
	}
	return nil
}
