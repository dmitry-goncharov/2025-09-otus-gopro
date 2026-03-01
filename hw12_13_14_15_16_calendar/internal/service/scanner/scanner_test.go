package scanner

import (
	"context"
	"testing"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorageScanner(t *testing.T) {
	t.Run("Clean events", func(t *testing.T) {
		ctx := context.Background()
		log := logger.NewMock()
		memStorage := memorystorage.New()
		storageScanner := NewStorageScanner(log, memStorage, 10*time.Minute)

		date := time.Now()
		y, m, d := date.Date()
		evt := storage.Event{
			ID:     uuid.New().String(),
			Title:  "Event",
			Date:   time.Date(y, m, d, date.Hour(), date.Minute()+5, 0, 0, date.Location()),
			UserID: uuid.New().String(),
		}
		memStorage.CreateEvent(ctx, evt)

		evts := storageScanner.Scan(ctx)

		require.Len(t, evts, 1)
	})
}
