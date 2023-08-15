package core

import (
	"github.com/fushiliang321/go-core/config/init"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/exception"
	"sync"
	"sync/atomic"
)

type (
	Server interface {
		Start(wg *sync.WaitGroup)
	}
	serverStartMonitor struct {
		isFinish  atomic.Bool //是否完成启动
		awaitChan chan byte   //等待通道
	}
)

var (
	startOnce            sync.Once
	awaitStartFinishOnce sync.Once
	_serverStartMonitor  = &serverStartMonitor{
		awaitChan: make(chan byte),
	}
)

func init() {
	_serverStartMonitor.isFinish.Store(false)
}

func Start() {
	defer func() {
		exception.Listener("core start", recover())
	}()
	startOnce.Do(func() {
		event.Dispatch(event.NewRegistered(event.BeforeServerStart, nil))
		wg := &sync.WaitGroup{}
		servers := init.Get()
		for _, ser := range servers {
			ser.Start(wg)
		}
		event.Dispatch(event.NewRegistered(event.AfterServerStart, nil))
		wg.Wait()
		event.Dispatch(event.NewRegistered(event.ServerEnd, nil))
	})
}

// 等待所有服务启动完成
func AwaitStartFinish(funs ...func()) {
	defer func() {
		for _, fun := range funs {
			go fun()
		}
	}()
	if _serverStartMonitor.isFinish.Load() {
		//已经完成启动
		return
	}
	defer func() {
		recover()
	}()
	awaitStartFinishOnce.Do(func() {
		event.Listener(event.Listen{
			EventNames: []string{event.AfterServerStart},
			Process: func(registered event.Registered) {
				//启动完成后把等待通道关闭
				_serverStartMonitor.isFinish.Store(true)
				close(_serverStartMonitor.awaitChan)
				for len(_serverStartMonitor.awaitChan) > 0 {
					<-_serverStartMonitor.awaitChan
				}
			},
		})
	})
	//还没启动完成就等待启动完成
	_serverStartMonitor.awaitChan <- 1
}
