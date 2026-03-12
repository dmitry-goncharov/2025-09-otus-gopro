package memorymq

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue"
)

type MemoryMQ struct {
	log   app.Logger
	queue []app.Message
}

func NewMemoryMQ(log app.Logger) queue.MessageQueue {
	return &MemoryMQ{
		log: log,
	}
}

func (q *MemoryMQ) Produce(_ context.Context, message app.Message) error {
	q.log.Debug("produce message to memorymq")
	q.queue = append(q.queue, message)
	return nil
}

func (q *MemoryMQ) Consume(ctx context.Context) (<-chan app.Message, error) {
	q.log.Debug("cosume message from memorymq")
	out := make(chan app.Message)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for len(q.queue) > 0 {
					msg := q.queue[0]
					out <- msg
					q.queue = q.queue[1:]
				}
			}
		}
	}()
	return out, nil
}

func (q *MemoryMQ) Close() error {
	q.log.Debug("close memorymq")
	return nil
}

func (q *MemoryMQ) Queue() []app.Message {
	return q.queue
}
