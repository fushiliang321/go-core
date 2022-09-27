package server

import (
	config "github.com/fushiliang321/go-core/config/rpc"
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

func init() {
	checkPort := helper.GetEnvDefault("SERVER_PORT_JSONRPC_HTTP", "9000")
	server, _ = go_jsonrpc.NewServer("http", "", checkPort)
	ip = helper.GetLocalIP()
	port, _ = strconv.Atoi(checkPort)
}

func (Service) Start(wg *sync.WaitGroup) {
	consulConfig := config.Get()
	if len(consulConfig.Services) > 0 {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 启动服务监听
			server.Start()
		}(wg)

		for _, s := range consulConfig.Services {
			RegisterServer(reflect.Indirect(reflect.ValueOf(s)).Type().Name()+"Service", s)
		}
	}

}
