package server

import (
	config "github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/helper"
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
	ip = helper.GetLocalIP(consulConfig.Address)
}

func (Service) Start(wg *sync.WaitGroup) {
	consulConfig := config.Get()
	if consulConfig.Services == nil || len(consulConfig.Services) == 0 {
		return
	}
	initialize()
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 启动服务监听
		server.Start()
	}(wg)
	for _, s := range consulConfig.Services {
		RegisterServer(reflect.Indirect(reflect.ValueOf(s)).Type().Name()+"Service", s)
	}
	server.Register(new(Health))
}
