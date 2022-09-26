package rpc

type Rpc struct {
	Consumers []string
	Services  []any
}

var consul = &Rpc{
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *Rpc) {
	consul = config
}

func Get() *Rpc {
	return consul
}
