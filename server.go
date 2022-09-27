package core

import (
	"github.com/fushiliang321/go-core/amqp"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/exception"
	grpc "github.com/fushiliang321/go-core/grpc/server"
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
	once sync.Once
)

var servers = []Server{
	amqp.Service{},
	consul.Service{},
	rpc.Service{},
	grpc.Service{},
	task.Service{},
	rateLimit.Service{},
	server.Service{},
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
	once.Do(func() {
		wg := &sync.WaitGroup{}
		for _, ser := range servers {
			ser.Start(wg)
		}
		wg.Wait()
	})
}
