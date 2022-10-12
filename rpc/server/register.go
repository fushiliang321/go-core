package server

import (
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/hashicorp/consul/api"
)

func RegisterServer(name string, s any) {
	b, _ := consul.RegisterServer(name, "jsonrpc-http", ip, port, &api.AgentServiceCheck{
		HTTP:   fmt.Sprintf("http://%s:%d/", ip, port),
		Method: "POST",
		Body:   "{\"id\":\"\",\"jsonrpc\":\"2.0\",\"method\":\"./\",\"params\":{}}",
	})
	if b {
		server.Register(s)
	}
}
