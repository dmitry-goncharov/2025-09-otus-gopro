package producer

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Producer interface {
	Produce(ctx context.Context, message app.Message) error
	Close() error
}
