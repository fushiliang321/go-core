package jsonRpcHttp

type JsonRpcHttp struct {
	Host      string
	Port      int
	Consumers []string
	Services  []any
}

var consul = &JsonRpcHttp{
	Host:      "",
	Port:      0,
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *JsonRpcHttp) {
	consul = config
}

func Get() *JsonRpcHttp {
	return consul
}
