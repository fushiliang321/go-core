package jsonrpc

type JsonRpc struct {
	Host      string
	Port      int
	Consumers []string
	Services  []any
}

var consul = &JsonRpc{
	Host:      "",
	Port:      0,
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *JsonRpc) {
	consul = config
}

func Get() *JsonRpc {
	return consul
}
