package server

import (
	config "core/config/consul"
	"core/helper"
	go_jsonrpc "github.com/iloveswift/go-jsonrpc"
	"reflect"
	"sync"
)

type Service struct {
}

var server go_jsonrpc.ServerInterface
var checkPort string

func init() {
	checkPort = helper.GetEnvDefault("SERVER_PORT_JSONRPC_HTTP", "9000")
	server, _ = go_jsonrpc.NewServer("http", "", checkPort)
}

func (Service) Start(wg *sync.WaitGroup) {
	consulConfig := config.Get()
	if len(consulConfig.Services) > 0 {
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			// 启动服务监听
			server.Start()
			wg.Done()
		}(wg)

		for _, s := range consulConfig.Services {
			RegisterServer(reflect.Indirect(reflect.ValueOf(s)).Type().Name()+"Service", s)
		}
	}

}
