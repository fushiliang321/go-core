package consul

import (
	"gitee.com/zvc/go-core/config/consul"
	"github.com/hashicorp/consul/api"
	"sync"
)

var consulConfig *consul.Consul

func (Service) Start(wg *sync.WaitGroup) {
	consulConfig = consul.Get()
	initApiConfig()
	if len(consulConfig.Consumers) > 0 {
		go func() {
			// 获取服务信息
			getData()
		}()
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
