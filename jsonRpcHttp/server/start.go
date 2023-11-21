package server

import (
	config "github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/helper/system"
	"github.com/fushiliang321/jsonrpc"
	"reflect"
	"strconv"
	"sync"
)

type Service struct{}

var (
	server jsonrpc.ServerInterface
	ip     string
	port   int
)

func initialize() {
	jsonRpcHttpConfig := config.Get()
	port = jsonRpcHttpConfig.Port

	server = NewHttpServer(jsonRpcHttpConfig.Host, strconv.Itoa(port))

	consulConfig := consul.GetConfig()
	ip = system.GetLocalIP(consulConfig.Address)
}

func (*Service) Start(wg *sync.WaitGroup) {
	_config := config.Get()
	if _config.Services == nil || len(_config.Services) == 0 {
		return
	}
	serviceRegistrations = make(map[string]*registerInfo, len(_config.Services))
	initialize()
	wg.Add(1)
	event.Dispatch(event.NewRegistered(event.BeforeJsonRpcServerStart, nil))
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 启动服务监听
		server.Start()
	}(wg)
	for _, s := range _config.Services {
		RegisterServer(reflect.Indirect(reflect.ValueOf(s)).Type().Name()+"Service", s)
	}
	server.Register(new(Health))
	event.Dispatch(event.NewRegistered(event.AfterJsonRpcServerStart))
}

func (*Service) PreEvents() []string {
	return []string{event.AfterLoggerServerStart}
}
