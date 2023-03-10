package server

import (
	"errors"
	"github.com/fushiliang321/go-core/consul"
	"github.com/hashicorp/consul/api"
	"time"
)

type Health struct{}

type params struct {
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

var (
	resultSuccess        = "success"
	resultError          = "error"
	serviceRegistrations *map[string]*api.AgentServiceRegistration //全局的服务注册信息
)

func (s *Health) Check(params *params) (*string, error) {
	if serviceRegistrations == nil {
		serviceRegistrations = consul.ServiceRegistrations()
		if serviceRegistrations == nil {
			//需要延迟响应，等待客户端请求超时
			time.Sleep(time.Minute)
			return &resultError, errors.New("服务不存在")
		}
	}
	registration, ok := (*serviceRegistrations)[params.Name]
	if !ok {
		//需要延迟响应，等待客户端请求超时
		time.Sleep(time.Minute)
		return &resultError, errors.New("服务不存在")
	}
	if protocol := registration.Meta["Protocol"]; protocol != params.Protocol {
		//需要延迟响应，等待客户端请求超时
		time.Sleep(time.Minute)
		return &resultError, errors.New("服务不存在")
	}
	return &resultSuccess, nil
}
