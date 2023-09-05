package dispatch

import (
	"fmt"
	config "github.com/fushiliang321/go-core/config/middlewares"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/response"
	"github.com/fushiliang321/go-core/middleware"
	types2 "github.com/fushiliang321/go-core/router/types"
	"github.com/fushiliang321/go-core/server/types"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"strings"
)

const gzipMinSize = 10000 //触发gzip压缩的最小长度

// Dispatch 包装整合路由调度的中间件
func Dispatch(handler types2.RequestHandler) fasthttp.RequestHandler {
	coreMiddlewares := middleware.GetCoreMiddlewares(handler)
	middlewares := config.Get()
	middlewaresHttp := append(middlewares.Http, coreMiddlewares.Http...)
	middlewaresWS := append(middlewares.WS, coreMiddlewares.WS...)
	middlewaresHttpLen := len(middlewaresHttp)
	middlewaresWsLen := len(middlewaresWS)

	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				response.ErrorResponse((*types2.RequestCtx)(ctx), 500, fmt.Sprintln("server exception:", err), nil)
				logger.Error("server exception:", fmt.Sprint(err))
				exception.Listener("server exception", err)
			}
		}()

		var (
			handlers  requestHandler
			_ctx      *types2.RequestCtx
			_, httpOk = ctx.UserValue(types.SERVER_HTTP_KEY).(*server.Server)
			_, wsOk   = ctx.UserValue(types.SERVER_WEBSOCKET_KEY).(*server.Server)
		)
		switch {
		case wsOk && "websocket" == strings.ToLower(strconv.B2S(ctx.Request.Header.Peek("upgrade"))):
			//优先判断websocket请求，避免同时开启http和websocket时无法升级到websocket
			handlers = requestHandler{
				len:         middlewaresWsLen,
				middlewares: make([]config.Middleware, middlewaresWsLen),
			}
			copy(handlers.middlewares, middlewaresWS)
		case httpOk:
			handlers = requestHandler{
				len:         middlewaresHttpLen,
				middlewares: make([]config.Middleware, middlewaresHttpLen),
			}
			copy(handlers.middlewares, middlewaresHttp)
			ctx.RemoveUserValue(types.SERVER_WEBSOCKET_KEY)
		default:
			//未知的协议，就直接返回空数据
			ctx.Write([]byte{})
			return
		}
		ctx.RemoveUserValue(types.SERVER_HTTP_KEY)
		_ctx = (*types2.RequestCtx)(ctx)
		_ctx.WriteAny(handlers.Process(_ctx))
	}
}
