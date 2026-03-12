package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
	// Posgresql driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	dsn *string
	log app.Logger
}

func New(dsn string, log app.Logger) app.Storage {
	return &Storage{
		dsn: &dsn,
		log: log,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	s.log.Debug("connect to sql storage")

	db, err := sqlx.Open("pgx", *s.dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	s.db = db
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.log.Debug("close sql storage")

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, evt storage.Event) error {
	s.log.Debug("create event in sql storage, evt:" + evt.String())

	query := `INSERT INTO events (id, title, date, user_id) VALUES (:id, :title, :date, :user_id);`

	res, err := s.db.NamedExecContext(ctx, query, &evt)
	if err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if ra == 0 {
		return storage.ErrAlreadyExists
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, evtID string, evt storage.Event) error {
	s.log.Debug("update event in sql storage, evtID:" + evtID + ", evt:" + evt.String())

	query := `UPDATE events SET title = :title WHERE id = :id;`
	args := map[string]any{
		"title": evt.Title,
		"id":    evtID,
	}

	res, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if ra == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, evtID string) error {
	s.log.Debug("delete event in sql storage, evtID:" + evtID)

	query := `DELETE FROM events WHERE id = :id;`
	args := map[string]any{
		"id": evtID,
	}

	res, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if ra == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) GetDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.log.Debug("get day events from sql storage, date:" + date.String())

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 1)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.log.Debug("get week events from sql storage, date:" + date.String())

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 7)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.log.Debug("get month events from sql storage, date:" + date.String())

	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 1, 0)

	return s.GetEventsByRange(ctx, begin, end)
}

func (s *Storage) GetEventsByRange(ctx context.Context, begin time.Time, end time.Time) ([]storage.Event, error) {
	query := `select * from events where date >= :begin and date < :end;`
	args := map[string]any{
		"begin": begin,
		"end":   end,
	}

	rows, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	events := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("error scaning query result: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) DeleteOutdatedEvents(ctx context.Context, date time.Time) error {
	query := `delete from events where date < :date`
	args := map[string]any{
		"date": date,
	}

	_, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}
