package core

import (
	"core/amqp"
	"core/consul"
	"core/exception"
	"core/rateLimit"
	rpc "core/rpc/server"
	"core/server"
	"core/task"
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
