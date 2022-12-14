package middleware

import (
	"fmt"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/router/types"
	"github.com/valyala/fasthttp"
)

type HttpCoreMiddleware struct {
	Handler types.RequestHandler
}

func (m *HttpCoreMiddleware) Process(ctx *fasthttp.RequestCtx, handler types.RequestHandler) (res any) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Response.SetStatusCode(500)
			res = helper.Error(500, fmt.Sprintln("server process exception:", err), nil)
			exception.Listener("server process exception", err)
		}
	}()
	res = m.Handler(ctx)
	return
}
