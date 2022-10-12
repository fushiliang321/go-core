package server

import (
	"encoding/json"
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/hashicorp/consul/api"
)

type checkBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  params `json:"params"`
}

func RegisterServer(name string, s any) {
	body := checkBody{
		Jsonrpc: "2.0",
		Method:  "Health.Check",
		Params: params{
			Name: name,
		},
	}
	bodyStr, _ := json.Marshal(body)
	b, _ := consul.RegisterServer(name, "jsonrpc-http", ip, port, &api.AgentServiceCheck{
		HTTP:   fmt.Sprintf("http://%s:%d/", ip, port),
		Method: "POST",
		Body:   string(bodyStr),
	})
	if b {
		server.Register(s)
	}
}
