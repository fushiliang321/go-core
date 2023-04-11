package server

import (
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/helper"
	"github.com/hashicorp/consul/api"
)

type (
	checkBody struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  *registerInfo `json:"params"`
	}
	registerInfo struct {
		Name     string `json:"name,omitempty"`
		Protocol string `json:"protocol,omitempty"`
		address  string
		port     int
	}
)

var serviceRegistrations map[string]*registerInfo //全局的服务注册信息

func RegisterServer(name string, s any) {
	var (
		_registerInfo = &registerInfo{
			Name:     name,
			Protocol: "jsonrpc-http",
			address:  ip,
			port:     port,
		}
		bodyStr, _ = helper.JsonEncode(checkBody{
			Jsonrpc: "2.0",
			Method:  "Health.Check",
			Params:  _registerInfo,
		})
	)
	b, _ := consul.RegisterServer(_registerInfo.Name, _registerInfo.Protocol, _registerInfo.address, _registerInfo.port, &api.AgentServiceCheck{
		HTTP:   fmt.Sprintf("http://%s:%d/", _registerInfo.address, _registerInfo.port),
		Method: "POST",
		Body:   bodyStr,
	})
	if b {
		server.Register(s)
		event.Dispatch(event.NewRegistered(event.JsonRpcServerRegister, *_registerInfo))
		serviceRegistrations[_registerInfo.Name] = _registerInfo
	}
}
