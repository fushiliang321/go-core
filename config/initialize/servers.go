package initialize

import (
	"github.com/fushiliang321/go-core/amqp"
	"github.com/fushiliang321/go-core/consul"
	grpc "github.com/fushiliang321/go-core/grpc/server"
	jsonRpcHttp "github.com/fushiliang321/go-core/jsonRpcHttp/server"
	"github.com/fushiliang321/go-core/logger"
	"github.com/fushiliang321/go-core/rateLimit"
	"github.com/fushiliang321/go-core/server"
	"github.com/fushiliang321/go-core/task"
	"sync"
)

type Service interface {
	Start(wg *sync.WaitGroup)
}

var servers = []Service{
	&logger.Service{},
	&amqp.Service{},
	&consul.Service{},
	&jsonRpcHttp.Service{},
	&grpc.Service{},
	&task.Service{},
	&rateLimit.Service{},
	&server.Service{},
}

func Register(s Service) {
	servers = append(servers, s)
}

func Registers(sers []Service) {
	servers = append(servers, sers...)
}

func Get() []Service {
	return servers
}
