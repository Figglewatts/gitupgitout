package ctxcron

import (
	"context"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	runner *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{runner: cron.New()}
}

func (s *CronScheduler) Run(ctx context.Context) {
	s.runner.Start()
	for range ctx.Done() {
		stopCtx := s.runner.Stop()
		<-stopCtx.Done()
		return
	}
}

func (s *CronScheduler) AddFunc(
	spec string, cmd func(),
) error {
	_, err := s.runner.AddFunc(spec, cmd)
	return err
}
