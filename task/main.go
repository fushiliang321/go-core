package task

import (
	"github.com/fushiliang321/go-core/config/task"
	"github.com/robfig/cron"
	"sync"
)

type Service struct{}

func (Service) Start(wg *sync.WaitGroup) {
	config := task.Get()
	if len(config.Crontabs) == 0 {
		return
	}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		c := cron.New()
		for _, crontab := range config.Crontabs {
			c.AddFunc(crontab.Rule, crontab.Callback)
		}
		c.Run()
	}(wg)
}
