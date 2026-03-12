package cleaner

import (
	"context"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Cleaner interface {
	Clean(ctx context.Context)
}

type StorageCleaner struct {
	log     app.Logger
	storage app.Storage
	date    time.Duration
}

func NewStorageCleaner(log app.Logger, storage app.Storage, date time.Duration) Cleaner {
	return &StorageCleaner{
		log:     log,
		storage: storage,
		date:    date,
	}
}

func (c *StorageCleaner) Clean(ctx context.Context) {
	date := time.Now().Add(-c.date)
	c.log.Debug("clean storage before " + date.String())
	c.storage.DeleteOutdatedEvents(ctx, date)
}
