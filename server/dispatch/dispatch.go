package dispatch

import (
	"fmt"
	config "github.com/fushiliang321/go-core/config/middlewares"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/middleware"
	types2 "github.com/fushiliang321/go-core/router/types"
	"github.com/fushiliang321/go-core/server/types"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

const gzipMinSize = 10000 //触发gzip压缩的最小长度

// 包装整合路由调度的中间件
func Dispatch(handler types2.RequestHandler) fasthttp.RequestHandler {
	coreMiddlewares := middleware.GetCoreMiddlewares(handler)
	middlewares := &config.Middlewares{
		Http: []config.Middleware{},
		WS:   []config.Middleware{},
	}
	*middlewares = *config.Get()
	middlewares.Http = append(middlewares.Http, coreMiddlewares.Http...)
	middlewares.WS = append(middlewares.WS, coreMiddlewares.WS...)
	middlewaresHttpLen := len(middlewares.Http)
	middlewaresWsLen := len(middlewares.WS)

	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				helper.ErrorResponse(ctx, 500, fmt.Sprintln("server exception:", err), nil)
				exception.Listener("server exception", err)
			}
		}()
		_, httpOk := ctx.UserValue(types.SERVER_HTTP_KEY).(*server.Server)
		_, wsOk := ctx.UserValue(types.SERVER_WEBSOCKET_KEY).(*server.Server)
		handlers := requestHandler{}
		switch {
		case wsOk && "websocket" == strings.ToLower(string(ctx.Request.Header.Peek("Upgrade"))):
			//优先判断websocket请求，避免同时开启http和websocket时无法升级到websocket
			handlers.len = middlewaresWsLen
			handlers.middlewares = make([]config.Middleware, handlers.len)
			copy(handlers.middlewares, middlewares.WS)
		case httpOk:
			handlers.len = middlewaresHttpLen
			handlers.middlewares = make([]config.Middleware, handlers.len)
			copy(handlers.middlewares, middlewares.Http)
			ctx.RemoveUserValue(types.SERVER_WEBSOCKET_KEY)
		default:
			//未知的协议，就直接返回空数据
			ctx.Write([]byte{})
			return
		}
		ctx.RemoveUserValue(types.SERVER_HTTP_KEY)
		res := handlers.Process(ctx)
		write(ctx, res)
	}
}

func write(ctx *fasthttp.RequestCtx, data any) {
	bytes, err := helper.AnyToBytes(data)
	if err != nil {
		log.Printf("server result err:%s\n", err)
		return
	}
	if len(*bytes) > gzipMinSize {
		ctx.Response.Header.Add("Content-Encoding", "gzip")
		ctx.Write(fasthttp.AppendGzipBytes([]byte{}, *bytes))
	} else {
		ctx.Write(*bytes)
	}
}
