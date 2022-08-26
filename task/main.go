package task

import (
	"core/config/task"
	"github.com/robfig/cron"
	"sync"
)

type Service struct {
}

func (Service) Start(wg *sync.WaitGroup) {
	config := task.Get()
	if len(config.Crontabs) == 0 {
		return
	}
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		c := cron.New()
		for _, crontab := range config.Crontabs {
			c.AddFunc(crontab.Rule, crontab.Callback)
		}
		c.Run()
	}(wg)
}
