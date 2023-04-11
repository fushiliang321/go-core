package core

import (
	"github.com/fushiliang321/go-core/amqp"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/exception"
	grpc "github.com/fushiliang321/go-core/grpc/server"
	jsonRpcHttp "github.com/fushiliang321/go-core/jsonRpcHttp/server"
	"github.com/fushiliang321/go-core/rateLimit"
	"github.com/fushiliang321/go-core/server"
	"github.com/fushiliang321/go-core/task"
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
	servers              = []Server{
		&amqp.Service{},
		&consul.Service{},
		&jsonRpcHttp.Service{},
		&grpc.Service{},
		&task.Service{},
		&rateLimit.Service{},
		&server.Service{},
	}
	_serverStartMonitor = &serverStartMonitor{
		awaitChan: make(chan byte),
	}
)

func init() {
	_serverStartMonitor.isFinish.Store(false)
}

func Register(s Server) {
	servers = append(servers, s)
}
func Registers(sers []Server) {
	servers = append(servers, sers...)
}

func Start() {
	defer func() {
		exception.Listener("core start", recover())
	}()
	startOnce.Do(func() {
		event.Dispatch(event.NewRegistered(event.BeforeServerStart, nil))
		wg := &sync.WaitGroup{}
		for _, ser := range servers {
			ser.Start(wg)
		}
		event.Dispatch(event.NewRegistered(event.AfterServerStart, nil))
		wg.Wait()
		event.Dispatch(event.NewRegistered(event.ServerEnd, nil))
	})
}

// 等待所有服务启动完成
func AwaitStartFinish() {
	if _serverStartMonitor.isFinish.Load() {
		//已经完成启动
		return
	}
	awaitStartFinishOnce.Do(func() {
		event.Listener(event.Listen{
			EventNames: []string{event.AfterServerStart},
			Process: func(registered event.Registered) {
				//启动完成后把等待通道关闭
				for range _serverStartMonitor.awaitChan {
				}
				close(_serverStartMonitor.awaitChan)
			},
		})
	})
	//还没启动完成就等待启动完成
	_serverStartMonitor.awaitChan <- 1
}
