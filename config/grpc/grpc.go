package grpc

type Service struct {
	RegisterFun any
	Handle      any
}

type Grpc struct {
	Host      string
	Port      int
	Services  []Service
	Consumers []string
}

var data = &Grpc{
	Services:  []Service{},
	Consumers: []string{},
}

func Set(config *Grpc) {
	data = config
}

func Get() *Grpc {
	return data
}
