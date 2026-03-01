package queue

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type MessageQueue interface {
	Produce(ctx context.Context, message app.Message) error
	Consume(ctx context.Context) (<-chan app.Message, error)
	Close() error
}
