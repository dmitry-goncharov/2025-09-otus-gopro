package consumer

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Cosumer interface {
	Consume(ctx context.Context) (<-chan app.Message, error)
	Close() error
}
