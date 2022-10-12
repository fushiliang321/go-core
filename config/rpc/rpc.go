package rpc

type Rpc struct {
	Host      string
	Port      int
	Consumers []string
	Services  []any
}

var consul = &Rpc{
	Host:      "",
	Port:      0,
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *Rpc) {
	consul = config
}

func Get() *Rpc {
	return consul
}
