package server

import (
	"github.com/fushiliang321/fasthttp2"
	"github.com/fushiliang321/go-core/config/routers"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/server/types"
	"github.com/fushiliang321/go-core/server/websocket"
	"github.com/valyala/fasthttp"
	"sync"
)

type Service struct{}

func (*Service) Start(wg *sync.WaitGroup) {
	var (
		config          = server.Get()
		serverConfigMap = map[string]map[byte]server.Server{}
		startWs         = false
	)

	//整理出各个端口需要开启的协议的配置
	for i := range config.Servers {
		var (
			ser  = config.Servers[i]
			addr = ser.Host + ":" + ser.Port
		)
		if _, ok := serverConfigMap[addr]; !ok {
			serverConfigMap[addr] = make(map[byte]server.Server, 2)
		}
		switch ser.Type {
		case types.SERVER_WEBSOCKET:
			startWs = true
		case types.SERVER_HTTP2:
			//如果存在http2，则不启动http1
			delete(serverConfigMap[addr], types.SERVER_HTTP)
		case types.SERVER_HTTP:
			//如果存在http2，则不启动http1
			if _, exist := serverConfigMap[addr][types.SERVER_HTTP2]; exist {
				continue
			}
		}
		serverConfigMap[addr][ser.Type] = ser
	}

	if len(serverConfigMap) == 0 {
		return
	}

	if startWs {
		websocket.Start()
	}

	for addr, serTypeMap := range serverConfigMap {
		var (
			httpServer *server.Server
			wsServer   *server.Server
		)
		for _type := range serTypeMap {
			var (
				_ser      = serTypeMap[_type]
				TLSConfig = mergeTLSConfig(config.Settings.TLS, _ser.TLS)
			)
			switch _ser.Type {
			case types.SERVER_WEBSOCKET:
				if _ser.Server == nil {
					wsServer = &_ser
				} else {
					listenAndServe(wg, _ser.Server, nil, &_ser, addr, TLSConfig)
				}
			case types.SERVER_HTTP:
				if _ser.Server == nil {
					httpServer = &_ser
				} else {
					listenAndServe(wg, _ser.Server, &_ser, nil, addr, TLSConfig)
				}
			case types.SERVER_HTTP2:
				fasthttp2.ConfigureServer(_ser.Server, fasthttp2.ServerConfig{})
				if _ser.Server == nil {
					httpServer = &_ser
				} else {
					listenAndServe(wg, _ser.Server, &_ser, nil, addr, TLSConfig)
				}
			}
		}
		if httpServer != nil || wsServer != nil {
			var (
				httpTls *server.TLS
				wsTls   *server.TLS
			)
			if httpServer != nil {
				httpTls = httpServer.TLS
			}
			if wsServer != nil {
				wsTls = wsServer.TLS
			}
			TLSConfig := mergeTLSConfig(config.Settings.TLS, wsTls, httpTls) //合并tls配置，websocket和http同时配置了tls优先使用http的tls配置
			listenAndServe(wg, &fasthttp.Server{}, httpServer, wsServer, addr, TLSConfig)
		}
	}
}

func mergeTLSConfig(configs ...*server.TLS) (config *server.TLS) {
	for _, c := range configs {
		if c == nil || c.KeyFile == "" || c.CertFile == "" {
			continue
		}
		if config == nil {
			config = &server.TLS{
				KeyFile:  c.KeyFile,
				CertFile: c.CertFile,
			}
		} else {
			config.KeyFile = c.KeyFile
			config.CertFile = c.CertFile
		}
	}
	return
}

func (*Service) PreEvents() []string {
	return []string{event.AfterLoggerServerStart}
}

func listenAndServeCommon(wg *sync.WaitGroup, serve *fasthttp.Server, httpServer, wsServer *server.Server, addr string, fun func() error) {
	serve.Handler = generateHandler(httpServer, wsServer, addr)
	if serve.Handler == nil {
		return
	}
	wg.Add(1)
	go func(isHttpServer, isWsServer bool, addr string) {
		if err := fun(); err != nil {
			logger.Warn("start fasthttp fail", err.Error())
		}
		serverEnd(isHttpServer, isWsServer, addr)
		wg.Done()
	}(httpServer != nil, wsServer != nil, addr)
}

func listenAndServe(wg *sync.WaitGroup, serve *fasthttp.Server, httpServer, wsServer *server.Server, addr string, TLSConfig *server.TLS) {
	if TLSConfig == nil {
		listenAndServeCommon(wg, serve, httpServer, wsServer, addr, func() error {
			return serve.ListenAndServe(addr)
		})
	} else {
		//配置了tls自动升级为http2
		fasthttp2.ConfigureServer(serve, fasthttp2.ServerConfig{})
		listenAndServeCommon(wg, serve, httpServer, wsServer, addr, func() error {
			return serve.ListenAndServeTLS(addr, TLSConfig.CertFile, TLSConfig.KeyFile)
		})
	}
}

func generateHandler(httpServer, wsServer *server.Server, addr string) func(ctx *fasthttp.RequestCtx) {
	routerConfig := routers.Get()
	if httpServer != nil {
		event.Dispatch(event.NewRegistered(event.HttpServerListen, addr))
		if wsServer == nil {
			//http服务器
			return func(ctx *fasthttp.RequestCtx) {
				ctx.SetUserValue(types.SERVER_HTTP_KEY, true)
				routerConfig.Handler(ctx)
			}
		} else {
			//http+ws服务器
			event.Dispatch(event.NewRegistered(event.WebsocketServerListen, addr))
			return func(ctx *fasthttp.RequestCtx) {
				ctx.SetUserValue(types.SERVER_HTTP_KEY, true)
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

func serverEnd(isHttpServer, isWsServer bool, addr string) {
	switch {
	case isHttpServer:
		event.Dispatch(event.NewRegistered(event.HttpServerListenEnd, addr))
	case isWsServer:
		event.Dispatch(event.NewRegistered(event.WebsocketServerListenEnd, addr))
	}
}
