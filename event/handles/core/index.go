package core

import (
	"github.com/fushiliang321/go-core/event"
	"sync"
)

// AwaitStartFinish 等待所有服务启动完成
func AwaitStartFinish(funs ...func()) {
	defer func() {
		for _, fun := range funs {
			go fun()
		}
	}()
	AwaitEvents([]string{event.AfterServerStart})
}

// 等待指定核心事件全部触发
func AwaitEvents(eventNames []string) {
	wg := &sync.WaitGroup{}
	func() {
		event.CoreServiceEventLog.RLock()
		defer event.CoreServiceEventLog.RUnlock()
		for _, name := range eventNames {
			if event.CoreServiceEventLog.AlreadyTriggered(name) {
				continue
			}

			var listen event.Listen
			listen = event.Listen{
				EventNames: []string{name},
				Process: func(registered event.Registered) {
					wg.Done()
					event.RemoveListener(&listen)
				},
			}
			wg.Add(1)
			event.AddListener(&listen)
		}
	}()
	wg.Wait()
}
