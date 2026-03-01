package sender

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Sender interface {
	Send(ctx context.Context, message <-chan app.Message)
}

type LogSender struct {
	log app.Logger
}

func NewLogSender(log app.Logger) Sender {
	return &LogSender{
		log: log,
	}
}

func (s *LogSender) Send(ctx context.Context, messages <-chan app.Message) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messages:
				s.log.Debug("send message " + message.String())
			}
		}
	}()
}
