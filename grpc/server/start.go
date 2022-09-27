package server

import (
	"github.com/fushiliang321/go-core/config/grpc"
	"sync"
)

type Service struct {
}

var config *grpc.Grpc

func (Service) Start(wg *sync.WaitGroup) {
	config = grpc.Get()
	if config.Services != nil && len(config.Services) > 0 {
		server := listen(config.Host, config.Port)
		for fun, srv := range config.Services {
			server.RegisterServer(fun, srv)
		}
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			// 启动服务监听
			server.Serve()
			wg.Done()
		}(wg)
	}
}
