package grpc

type (
	Service struct {
		RegisterFun any
		Handle      any
	}
	Grpc struct {
		Host                   string
		Port                   int
		ConnectMaxMultiplexNum uint32 //连接最大复用次数
		Services               []Service
		Consumers              []string
	}
)

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
