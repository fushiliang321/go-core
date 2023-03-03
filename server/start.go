package server

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/routers"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/server/types"
	"github.com/fushiliang321/go-core/server/websocket"
	"github.com/valyala/fasthttp"
	"sync"
)

type Service struct {
}

func (Service) Start(wg *sync.WaitGroup) {
	r := routers.Get()
	config := server.Get()

	serverMap := map[string]map[byte]server.Server{}

	startWs := false
	for _, s := range config.Servers {
		addr := s.Host + ":" + s.Port
		if _, ok := serverMap[addr]; !ok {
			serverMap[addr] = make(map[byte]server.Server, 2)
		}
		serverMap[addr][s.Type] = s
		if s.Type == types.SERVER_WEBSOCKET {
			startWs = true
		}
	}

	if startWs {
		websocket.Start()
	}

	for addr, s := range serverMap {
		wg.Add(1)
		go func(addr string, sers map[byte]server.Server, wg *sync.WaitGroup) {
			defer wg.Done()
			var (
				httpServer *server.Server
				wsServer   *server.Server
				err        error
			)
			for _, ser := range sers {
				switch ser.Type {
				case types.SERVER_WEBSOCKET:
					wsServer = &ser
				case types.SERVER_HTTP:
					httpServer = &ser
				}
			}

			if httpServer != nil {
				if wsServer == nil {
					//http服务器
					err = fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
						ctx.SetUserValue(types.SERVER_HTTP_KEY, httpServer)
						r.Handler(ctx)
					})
				} else {
					//http+ws服务器
					err = fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
						ctx.SetUserValue(types.SERVER_HTTP_KEY, httpServer)
						ctx.SetUserValue(types.SERVER_WEBSOCKET_KEY, wsServer)
						r.Handler(ctx)
					})
				}
			} else if wsServer != nil {
				//ws服务器
				err = fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
					ctx.SetUserValue(types.SERVER_WEBSOCKET_KEY, wsServer)
					r.Handler(ctx)
				})
			}

			if err != nil {
				fmt.Println("start fasthttp fail", err.Error())
			}
		}(addr, s, wg)
	}
}
