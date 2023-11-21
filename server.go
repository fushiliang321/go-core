package core

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/initialize"
	"github.com/fushiliang321/go-core/config/initialize/service"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/event/handles/core"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"reflect"
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
		var (
			servers    = initialize.Get()
			serStartWg = &sync.WaitGroup{} //等待所有服务启动
			serEndWg   = &sync.WaitGroup{} //等待所有服务结束
		)
		for _, ser := range servers {
			serStartWg.Add(1)
			go func(ser service.Service, serEndWg, serStartWg *sync.WaitGroup) {
				serviceStart(ser, serEndWg)
				serStartWg.Done()
			}(ser, serEndWg, serStartWg)
		}
		serStartWg.Wait()
		event.Dispatch(event.NewRegistered(event.AfterServerStart))
		serEndWg.Wait()
		event.Dispatch(event.NewRegistered(event.ServerEnd))
	})
}

func serviceStart(ser service.Service, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("core start error:", reflect.ValueOf(ser).Elem().Type().String(), fmt.Sprint(err))
			exception.Listener("core start error", err)
		}
	}()
	preEvents := ser.PreEvents()
	if preEvents != nil {
		//等待前置事件触发
		core.AwaitEvents(preEvents)
	}
	ser.Start(wg)
}
