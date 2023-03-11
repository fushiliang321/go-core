package server

import (
	"errors"
	"time"
)

type Health struct{}

var (
	resultSuccess = "success"
	resultError   = "error"
)

func (s *Health) Check(params *registerInfo) (*string, error) {
	if len(serviceRegistrations) == 0 {
		//需要延迟响应，等待客户端请求超时
		time.Sleep(time.Minute)
		return &resultError, errors.New("服务不存在")
	}
	registration, ok := serviceRegistrations[params.Name]
	if !ok {
		//需要延迟响应，等待客户端请求超时
		time.Sleep(time.Minute)
		return &resultError, errors.New("服务不存在")
	}
	if registration.Protocol != params.Protocol {
		//需要延迟响应，等待客户端请求超时
		time.Sleep(time.Minute)
		return &resultError, errors.New("服务不存在")
	}
	return &resultSuccess, nil
}
