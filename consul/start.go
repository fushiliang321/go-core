package consul

import (
	"github.com/fushiliang321/go-core/config/consul"
	"github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/hashicorp/consul/api"
	"sync"
)

var consulConfig *consul.Consul

func (Service) Start(_ *sync.WaitGroup) {
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
		AddServices(consumerServiceNames)
	}
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
