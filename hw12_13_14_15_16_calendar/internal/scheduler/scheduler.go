package scheduler

import (
	"context"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
)

type Task func(ctx context.Context)

type Scheduler interface {
	Run(ctx context.Context)
}

type SimpleScheduler struct {
	name     string
	log      app.Logger
	interval time.Duration
	task     Task
}

func NewSimpleScheduler(name string, log app.Logger, interval time.Duration, task Task) Scheduler {
	return &SimpleScheduler{
		name:     name,
		log:      log,
		interval: interval,
		task:     task,
	}
}

func (s *SimpleScheduler) Run(ctx context.Context) {
	s.log.Info("run " + s.name)
	tick := time.NewTicker(s.interval)
	defer func() {
		s.log.Info("stop " + s.name)
		tick.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				s.log.Debug("run task of " + s.name)
				s.task(ctx)
			}
		}
	}
}
