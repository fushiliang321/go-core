package server

import (
	"github.com/fushiliang321/go-core/config/routers"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/logger"
	"github.com/fushiliang321/go-core/server/types"
	"github.com/fushiliang321/go-core/server/websocket"
	"github.com/valyala/fasthttp"
	"sync"
)

type Service struct{}

func (*Service) Start(wg *sync.WaitGroup) {
	var (
		config = server.Get()

		serverMap = map[string]map[byte]server.Server{}

		startWs = false

		ok bool
	)
	for i := range config.Servers {
		var (
			ser  = config.Servers[i]
			addr = ser.Host + ":" + ser.Port
		)
		if _, ok = serverMap[addr]; !ok {
			serverMap[addr] = make(map[byte]server.Server, 2)
		}
		serverMap[addr][ser.Type] = ser
		if ser.Type == types.SERVER_WEBSOCKET {
			startWs = true
		}
	}

	if len(serverMap) == 0 {
		return
	}

	if startWs {
		websocket.Start()
	}

	for addr, sers := range serverMap {
		var (
			httpServer *server.Server
			wsServer   *server.Server
		)
		for _, ser := range sers {
			switch ser.Type {
			case types.SERVER_WEBSOCKET:
				if ser.Server == nil {
					wsServer = &ser
				} else {
					listenAndServe(wg, ser.Server, nil, &ser, addr)
				}
			case types.SERVER_HTTP:
				if ser.Server == nil {
					httpServer = &ser
				} else {
					listenAndServe(wg, ser.Server, &ser, nil, addr)
				}
			}
		}
		if httpServer != nil || wsServer != nil {
			listenAndServe(wg, &fasthttp.Server{}, httpServer, wsServer, addr)
		}
	}
}

func listenAndServe(wg *sync.WaitGroup, serve *fasthttp.Server, httpServer, wsServer *server.Server, addr string) {
	serve.Handler = generateHandler(httpServer, wsServer, addr)
	if serve.Handler == nil {
		return
	}
	wg.Add(1)
	go func(addr string) {
		if err := serve.ListenAndServe(addr); err != nil {
			logger.Warn("start fasthttp fail", err.Error())
		}
		wg.Done()
	}(addr)
}

func generateHandler(httpServer, wsServer *server.Server, addr string) func(ctx *fasthttp.RequestCtx) {
	routerConfig := routers.Get()
	if httpServer != nil {
		if wsServer == nil {
			//http服务器
			event.Dispatch(event.NewRegistered(event.HttpServerListen, addr))
			return func(ctx *fasthttp.RequestCtx) {
				ctx.SetUserValue(types.SERVER_HTTP_KEY, httpServer)
				routerConfig.Handler(ctx)
			}
		} else {
			//http+ws服务器
			event.Dispatch(event.NewRegistered(event.HttpServerListen, addr))
			event.Dispatch(event.NewRegistered(event.WebsocketServerListen, addr))
			return func(ctx *fasthttp.RequestCtx) {
				ctx.SetUserValue(types.SERVER_HTTP_KEY, httpServer)
				ctx.SetUserValue(types.SERVER_WEBSOCKET_KEY, wsServer)
				routerConfig.Handler(ctx)
			}
		}
	} else if wsServer != nil {
		//ws服务器
		event.Dispatch(event.NewRegistered(event.WebsocketServerListen, addr))
		return func(ctx *fasthttp.RequestCtx) {
			ctx.SetUserValue(types.SERVER_WEBSOCKET_KEY, wsServer)
			routerConfig.Handler(ctx)
		}
	}
	return nil
}
