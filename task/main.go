package task

import (
	"github.com/fushiliang321/go-core/config/task"
	"github.com/fushiliang321/go-core/event"
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
		defer func() {
			recover()
			wg.Done()
		}()
		cron := &Cron{
			before: func(sign any) {
				event.NewRegistered(event.BeforeTaskExecute, sign)
			},
			after: func(sign any) {
				event.NewRegistered(event.AfterTaskExecute, sign)
			},
		}
		event.Dispatch(event.NewRegistered(event.BeforeTaskServerStart))
		for i := range config.Crontabs {
			crontab := config.Crontabs[i]
			callback := crontab.Callback
			err := cron.AddFuncSign(crontab.Rule, callback, crontab)
			if err != nil {
				continue
			}
			event.Dispatch(event.NewRegistered(event.TaskRegister, crontab))
		}
		event.Dispatch(event.NewRegistered(event.AfterTaskServerStart))
		cron.Run()
	}(wg)
}
