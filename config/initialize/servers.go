package initialize

import (
	"github.com/fushiliang321/go-core/amqp"
	"github.com/fushiliang321/go-core/config/initialize/Service"
	"github.com/fushiliang321/go-core/consul"
	grpc "github.com/fushiliang321/go-core/grpc/server"
	jsonRpcHttp "github.com/fushiliang321/go-core/jsonRpcHttp/server"
	"github.com/fushiliang321/go-core/logger"
	"github.com/fushiliang321/go-core/rateLimit"
	"github.com/fushiliang321/go-core/server"
	"github.com/fushiliang321/go-core/task"
)

var services = []service.Service{
	&logger.Service{},
	&amqp.Service{},
	&consul.Service{},
	&jsonRpcHttp.Service{},
	&grpc.Service{},
	&task.Service{},
	&rateLimit.Service{},
	&server.Service{},
}

func Register(s service.Service) {
	services = append(services, s)
}

func Registers(sers []service.Service) {
	services = append(services, sers...)
}

func Get() []service.Service {
	return services
}
