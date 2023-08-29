package middleware

import (
	"fmt"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/response"
	"github.com/fushiliang321/go-core/router/types"
)

type HttpCoreMiddleware struct {
	Handler types.RequestHandler
}

func (m *HttpCoreMiddleware) Process(ctx *types.RequestCtx, handler types.RequestHandler) (res any) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Response.SetStatusCode(500)
			res = response.Error(500, fmt.Sprintln("server process exception:", err), nil)
			logger.Error("server process exception:", err)
			exception.Listener("server process exception", err)
		}
	}()
	res = m.Handler(ctx)
	return
}
