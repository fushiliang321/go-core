package consul

import (
	"github.com/fushiliang321/go-core/config/consul"
	"github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/config/rpc"
	"github.com/hashicorp/consul/api"
	"sync"
)

var consulConfig *consul.Consul

func (Service) Start(wg *sync.WaitGroup) {
	consulConfig = consul.Get()
	rpcConfig := rpc.Get()
	grpcConfig := grpc.Get()
	initApiConfig()

	consumerServiceNames := []string{}
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
		getServiceData(consumerServiceNames)
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
