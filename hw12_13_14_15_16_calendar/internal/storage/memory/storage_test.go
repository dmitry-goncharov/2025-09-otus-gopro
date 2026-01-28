package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Create event", func(t *testing.T) {
		evt := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		err := s.CreateEvent(ctx, evt)

		require.NoError(t, err)
		require.Len(t, s.(*Storage).evts, 1)

		err = s.CreateEvent(ctx, evt)

		require.ErrorIs(t, err, storage.ErrAlreadyExists)
		require.Len(t, s.(*Storage).evts, 1)
	})

	t.Run("Update event", func(t *testing.T) {
		evt := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		err := s.CreateEvent(ctx, evt)

		require.NoError(t, err)
		require.Len(t, s.(*Storage).evts, 1)

		evt.Title = "Updated Event"

		err = s.UpdateEvent(ctx, evt.ID, evt)

		require.NoError(t, err)
		require.Len(t, s.(*Storage).evts, 1)
		require.Equal(t, "Updated Event", s.(*Storage).evts[evt.ID].Title)

		err = s.UpdateEvent(ctx, "id", evt)

		require.ErrorIs(t, err, storage.ErrNotFound)
		require.Len(t, s.(*Storage).evts, 1)
	})

	t.Run("Delete event", func(t *testing.T) {
		evt := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		err := s.CreateEvent(ctx, evt)

		require.NoError(t, err)
		require.Len(t, s.(*Storage).evts, 1)

		err = s.DeleteEvent(ctx, "id")

		require.ErrorIs(t, err, storage.ErrNotFound)
		require.Len(t, s.(*Storage).evts, 1)

		err = s.DeleteEvent(ctx, evt.ID)

		require.NoError(t, err)
		require.Len(t, s.(*Storage).evts, 0)
	})
}

func TestStorageGetEvents(t *testing.T) {
	t.Run("Get day events", func(t *testing.T) {
		date := time.Now()
		y, m, d := date.Date()
		evt1 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Date(y, m, d, 0, 0, 0, 0, date.Location()),
			UserID: uuid.New().String(),
		}
		evt2 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}
		evt3 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, -1),
			UserID: uuid.New().String(),
		}
		evt4 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, 1),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		createEvents(ctx, s, evt1, evt2, evt3, evt4)

		events, err := s.GetDayEvents(ctx, time.Now())

		require.NoError(t, err)
		require.Len(t, events, 2)
		orig := make([]storage.Event, 0, 2)
		orig = append(orig, evt1, evt2)
		require.ElementsMatch(t, orig, events)
	})

	t.Run("Get week events", func(t *testing.T) {
		evt1 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}
		evt2 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, -1),
			UserID: uuid.New().String(),
		}
		evt3 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, 6),
			UserID: uuid.New().String(),
		}
		evt4 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, 7),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		createEvents(ctx, s, evt1, evt2, evt3, evt4)

		events, err := s.GetWeekEvents(ctx, time.Now())

		require.NoError(t, err)
		require.Len(t, events, 2)
		orig := make([]storage.Event, 0, 2)
		orig = append(orig, evt1, evt3)
		require.ElementsMatch(t, orig, events)
	})

	t.Run("Get month events", func(t *testing.T) {
		evt1 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now(),
			UserID: uuid.New().String(),
		}
		evt2 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 0, -1),
			UserID: uuid.New().String(),
		}
		evt3 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 1, -1),
			UserID: uuid.New().String(),
		}
		evt4 := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Now().AddDate(0, 1, 0),
			UserID: uuid.New().String(),
		}

		ctx := context.Background()
		s := New()

		require.Len(t, s.(*Storage).evts, 0)

		createEvents(ctx, s, evt1, evt2, evt3, evt4)

		events, err := s.GetMonthEvents(ctx, time.Now())

		require.NoError(t, err)
		require.Len(t, events, 2)
		orig := make([]storage.Event, 0, 2)
		orig = append(orig, evt1, evt3)
		require.ElementsMatch(t, orig, events)
	})
}

func createEvents(ctx context.Context, s app.Storage, evts ...storage.Event) {
	for _, evt := range evts {
		s.CreateEvent(ctx, evt)
	}
}
