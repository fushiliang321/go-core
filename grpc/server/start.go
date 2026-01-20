package server

import (
	"sync"

	"github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/event"
)

type Service struct{}

var config *grpc.Grpc

func (*Service) Start(wg *sync.WaitGroup) {
	config = grpc.Get()
	if config.Services == nil || len(config.Services) == 0 {
		return
	}
	event.Dispatch(event.NewRegistered(event.BeforeGrpcServerStart))
	server := listen(config)
	regSuccess := false
	for _, service := range config.Services {
		if server.RegisterServer(service.Handle, service.RegisterFun) {
			regSuccess = true
		}
	}
	if regSuccess {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 启动服务监听
			server.Serve()
		}(wg)
	}
	event.Dispatch(event.NewRegistered(event.AfterGrpcServerStart))
}

func (*Service) PreEvents() []string {
	return []string{event.AfterLoggerServerStart}
}
