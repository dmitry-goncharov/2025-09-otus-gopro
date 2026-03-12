package queuefactory

import (
	"fmt"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue/memorymq"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue/rabbitmq"
)

const (
	IMqueue = "imq"
	RMQueue = "rmq"
)

func NewMessageQueue(log app.Logger, queueConf *config.QueueConf) (queue.MessageQueue, error) {
	switch queueConf.Type {
	case IMqueue:
		return memorymq.NewMemoryMQ(log), nil
	case RMQueue:
		return rabbitmq.NewRabbitMQ(log, queueConf.RMQ.Dsn, queueConf.RMQ.QName), nil
	default:
		return nil, fmt.Errorf("illegal message queue type: %s", queueConf.Type)
	}
}
