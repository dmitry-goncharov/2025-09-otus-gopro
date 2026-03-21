package sender

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Sender interface {
	Send(ctx context.Context, message <-chan app.Message)
}

type LogSender struct {
	log     app.Logger
	storage app.Storage
}

func NewLogSender(log app.Logger, storage app.Storage) Sender {
	return &LogSender{
		log:     log,
		storage: storage,
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
				err := s.storage.LogEventNotification(ctx, message.ID)
				if err != nil {
					s.log.Error("error logging event notification ID: " + message.ID)
				}
			}
		}
	}()
}
