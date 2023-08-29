package core

import (
	"github.com/fushiliang321/go-core/config/initialize"
	"github.com/fushiliang321/go-core/event"
	"sync"
)

var startOnce sync.Once

func Start() {
	defer func() {
		//if err := recover(); err != nil {
		//	logger.Error("core start error:", err)
		//	exception.Listener("core start", err)
		//}
	}()

	startOnce.Do(func() {
		event.Dispatch(event.NewRegistered(event.BeforeServerStart, nil))
		wg := &sync.WaitGroup{}
		servers := initialize.Get()
		for _, ser := range servers {
			ser.Start(wg)
		}
		event.Dispatch(event.NewRegistered(event.AfterServerStart, nil))
		wg.Wait()
		event.Dispatch(event.NewRegistered(event.ServerEnd, nil))
	})
}
