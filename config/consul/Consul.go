package consul

import (
	"github.com/hashicorp/consul/api"
)

type ConsulServiceHealthCheck struct {
	Timeout                        string //健康检测超时时间
	Interval                       string // 健康检查间隔
	DeregisterCriticalServiceAfter string // check失败后删除服务，注销时间，相当于过期时间
}

type Consul struct {
	api.Config
	HealthCheck *ConsulServiceHealthCheck
}

var consul = &Consul{
	HealthCheck: &ConsulServiceHealthCheck{},
}

func Set(config *Consul) {
	consul = config
}

func Get() *Consul {
	return consul
}
