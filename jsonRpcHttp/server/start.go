package server

import (
	config "github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/fushiliang321/go-core/helper"
	go_jsonrpc "github.com/iloveswift/go-jsonrpc"
	"reflect"
	"strconv"
	"sync"
)

type Service struct {
}

var server go_jsonrpc.ServerInterface
var ip string
var port int

func initialize() {
	consulConfig := config.Get()
	port = consulConfig.Port
	server, _ = go_jsonrpc.NewServer("http", consulConfig.Host, strconv.Itoa(port))
	ip = helper.GetLocalIP()
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
