package core

import (
	"github.com/fushiliang321/go-core/amqp/consumer"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/rateLimit"
	rpc "github.com/fushiliang321/go-core/rpc/server"
	"github.com/fushiliang321/go-core/server"
	"github.com/fushiliang321/go-core/task"
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
			consumer.Service{},
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
