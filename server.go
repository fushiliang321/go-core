package core

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/initialize"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"sync"
)

var startOnce sync.Once

func Start() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("core start error:", fmt.Sprint(err))
			exception.Listener("core start", err)
		}
	}()

	startOnce.Do(func() {
		event.Dispatch(event.NewRegistered(event.BeforeServerStart))
		wg := &sync.WaitGroup{}
		servers := initialize.Get()
		for _, ser := range servers {
			ser.Start(wg)
		}
		event.Dispatch(event.NewRegistered(event.AfterServerStart))
		wg.Wait()
		event.Dispatch(event.NewRegistered(event.ServerEnd))
	})
}
