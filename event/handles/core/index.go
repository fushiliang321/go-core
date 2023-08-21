package core

import (
	"github.com/fushiliang321/go-core/event"
	"sync"
	"sync/atomic"
)

type (
	serverStartMonitor struct {
		isFinish  atomic.Bool //是否完成启动
		awaitChan chan byte   //等待通道
	}
)

var (
	awaitStartFinishOnce sync.Once
	_serverStartMonitor  = &serverStartMonitor{
		awaitChan: make(chan byte),
	}
)

func init() {
	_serverStartMonitor.isFinish.Store(false)
	//监听服务启动完成事件
	event.Listener(event.Listen{
		EventNames: []string{event.AfterServerStart},
		Process: func(registered event.Registered) {
			//启动完成后把等待通道关闭
			_serverStartMonitor.isFinish.Store(true)
			close(_serverStartMonitor.awaitChan)
			for range _serverStartMonitor.awaitChan {
			}
		},
	})
}

// AwaitStartFinish 等待所有服务启动完成
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
				for range _serverStartMonitor.awaitChan {
				}
			},
		})
	})
	//还没启动完成就等待启动完成
	_serverStartMonitor.awaitChan <- 1
}
