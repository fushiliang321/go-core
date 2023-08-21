package core

import (
	"github.com/fushiliang321/go-core/config/initialize"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/exception"
	"sync"
)

var startOnce sync.Once

func Start() {
	defer func() {
		exception.Listener("core start", recover())
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
