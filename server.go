package core

import (
	"gitee.com/zvc/go-core/amqp"
	"gitee.com/zvc/go-core/consul"
	"gitee.com/zvc/go-core/exception"
	"gitee.com/zvc/go-core/rateLimit"
	rpc "gitee.com/zvc/go-core/rpc/server"
	"gitee.com/zvc/go-core/server"
	"gitee.com/zvc/go-core/task"
	"sync"
)

type Server interface {
	Start(wg *sync.WaitGroup)
}

var (
	servers []Server
	once    sync.Once
)

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
	once.Do(func() {
		wg := &sync.WaitGroup{}
		Registers([]Server{
			amqp.Service{},
			consul.Service{},
			rpc.Service{},
			task.Service{},
			rateLimit.Service{},
			server.Service{},
		})
		for _, ser := range servers {
			ser.Start(wg)
		}
		wg.Wait()
	})
}
