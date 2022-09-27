package grpc

import (
	"github.com/fushiliang321/go-core/helper"
	"strconv"
)

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
	if data.Port == 0 {
		data.Port, _ = strconv.Atoi(helper.GetEnvDefault("SERVER_PORT_GRPC", "9000"))
	}
	if data.Host == "" {
		data.Host = helper.GetEnvDefault("SERVER_ADDRESS_GRPC", "0.0.0.0")
	}
}

func Get() *Grpc {
	return data
}
