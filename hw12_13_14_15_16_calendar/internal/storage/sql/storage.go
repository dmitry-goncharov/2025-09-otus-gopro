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
}

func New(dsn string) app.Storage {
	return &Storage{
		dsn: &dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
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
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, evt storage.Event) error {
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
	query := `UPDATE events SET title = :title WHERE id = :id;`

	res, err := s.db.ExecContext(ctx, query, evt.Title, evtID)
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
	query := `DELETE FROM events WHERE id = :id;`

	res, err := s.db.ExecContext(ctx, query, evtID)
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
	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 1)

	return s.getEventsByRange(ctx, begin, end)
}

func (s *Storage) GetWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 0, 7)

	return s.getEventsByRange(ctx, begin, end)
}

func (s *Storage) GetMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	y, m, d := date.Date()
	begin := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	end := begin.AddDate(0, 1, 0)

	return s.getEventsByRange(ctx, begin, end)
}

func (s *Storage) getEventsByRange(ctx context.Context, begin time.Time, end time.Time) ([]storage.Event, error) {
	args := map[string]any{
		"begin": begin,
		"end":   end,
	}

	query := `select * from events where date >= :begin and date < :end;`

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
