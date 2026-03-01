package publisher

import (
	"context"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/service/producer"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/service/scanner"
)

type Publisher interface {
	Publish(ctx context.Context)
}

type ScaProPublisher struct {
	log      app.Logger
	scanner  scanner.Scanner
	producer producer.Producer
}

func NewScaProPublisher(log app.Logger, scanner scanner.Scanner, producer producer.Producer) Publisher {
	return &ScaProPublisher{
		log:      log,
		scanner:  scanner,
		producer: producer,
	}
}

func (p *ScaProPublisher) Publish(ctx context.Context) {
	p.log.Debug("publish")
	messages := p.scanner.Scan(ctx)
	for _, message := range messages {
		p.producer.Produce(ctx, message)
	}
}
