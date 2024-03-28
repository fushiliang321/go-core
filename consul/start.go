package consul

import (
	"github.com/fushiliang321/go-core/config/consul"
	"github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/config/initialize/service"
	"github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/fushiliang321/go-core/event"
	"github.com/hashicorp/consul/api"
	"sync"
)

var consulConfig *consul.Consul

type Service struct {
	service.BaseStruct
}

func (*Service) Start(_ *sync.WaitGroup) {
	consulConfig = consul.Get()
	rpcConfig := jsonRpcHttp.Get()
	grpcConfig := grpc.Get()
	initApiConfig()

	var consumerServiceNames []string
	if rpcConfig.Consumers != nil {
		for _, serviceName := range rpcConfig.Consumers {
			consumerServiceNames = append(consumerServiceNames, serviceName)
		}
	}
	if grpcConfig.Consumers != nil {
		for _, serviceName := range grpcConfig.Consumers {
			consumerServiceNames = append(consumerServiceNames, serviceName)
		}
	}
	if len(consumerServiceNames) > 0 {
		// 获取服务信息
		event.Dispatch(event.NewRegistered(event.BeforeConsulConsumerServerStart))
		AddServices(consumerServiceNames)
		event.Dispatch(event.NewRegistered(event.AfterConsulConsumerServerStart))
	}
	event.Dispatch(event.NewRegistered(event.ConsulInitFinish))
}

func (*Service) PreEvents() []string {
	return []string{event.AfterLoggerServerStart}
}

func initApiConfig() {
	apiConfig = api.DefaultConfig()
	if consulConfig.Address != "" {
		apiConfig.Address = consulConfig.Address
	}
}

func GetConfig() *api.Config {
	return apiConfig
}
