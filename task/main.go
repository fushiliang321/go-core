package task

import (
	"github.com/fushiliang321/go-core/config/task"
	"github.com/fushiliang321/go-core/event"
	"github.com/robfig/cron"
	"sync"
)

type Service struct{}

func (*Service) Start(wg *sync.WaitGroup) {
	config := task.Get()
	if len(config.Crontabs) == 0 {
		return
	}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		event.Dispatch(event.NewRegistered(event.BeforeTaskServerStart, nil))
		c := cron.New()
		for i := range config.Crontabs {
			crontab := config.Crontabs[i]
			BeforeTaskExecuteRegistered := event.NewRegistered(event.BeforeTaskExecute, crontab)
			AfterTaskExecuteRegistered := event.NewRegistered(event.AfterTaskExecute, crontab)
			err := c.AddFunc(crontab.Rule, func() {
				event.Dispatch(BeforeTaskExecuteRegistered)
				crontab.Callback()
				event.Dispatch(AfterTaskExecuteRegistered)
			})
			if err != nil {
				continue
			}
			event.Dispatch(event.NewRegistered(event.TaskRegister, crontab))
		}
		event.Dispatch(event.NewRegistered(event.AfterTaskServerStart, nil))
		c.Run()
	}(wg)
}
