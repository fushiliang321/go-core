package server

import (
	"github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/config/initialize/service"
	"github.com/fushiliang321/go-core/event"
	"sync"
)

type Service struct {
	service.BaseStruct
}

var config *grpc.Grpc

func (*Service) Start(wg *sync.WaitGroup) {
	config = grpc.Get()
	if config.Services == nil || len(config.Services) == 0 {
		return
	}
	event.Dispatch(event.NewRegistered(event.BeforeGrpcServerStart))
	server := listen(config.Host, config.Port)
	for _, service := range config.Services {
		server.RegisterServer(service.Handle, service.RegisterFun)
	}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 启动服务监听
		server.Serve()
	}(wg)
	event.Dispatch(event.NewRegistered(event.AfterGrpcServerStart))
}
