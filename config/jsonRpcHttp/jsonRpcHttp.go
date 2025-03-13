package jsonRpcHttp

import "golang.org/x/net/context"

type JsonRpcHttp struct {
	Host           string
	Port           int
	Consumers      []string
	Services       []any
	DefaultContext func() context.Context
}

var jsonRpcHttp = &JsonRpcHttp{
	Host:      "",
	Port:      0,
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *JsonRpcHttp) {
	jsonRpcHttp = config
}

func Get() *JsonRpcHttp {
	return jsonRpcHttp
}
