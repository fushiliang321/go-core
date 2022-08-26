package server

import (
	"core/config/routers"
	"core/config/server"
	"core/server/types"
	"core/server/websocket"
	"fmt"
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
			http := &server.Server{}
			ws := http
			for t := range sers {
				ser := sers[t]
				switch ser.Type {
				case types.SERVER_WEBSOCKET:
					ws = &ser
				case types.SERVER_HTTP:
					http = &ser
				}
			}
			if err := fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
				if ws.Type == types.SERVER_WEBSOCKET {
					ctx.SetUserValue(types.SERVER_WEBSOCKET_KEY, ws)
				}
				if http.Type == types.SERVER_HTTP {
					ctx.SetUserValue(types.SERVER_HTTP_KEY, http)
				}
				r.Handler(ctx)
			}); err != nil {
				fmt.Println("start fasthttp fail", err.Error())
			}
			wg.Done()
		}(addr, s, wg)
	}
}
