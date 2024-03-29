package dispatch

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/middlewares"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/response"
	"github.com/fushiliang321/go-core/router/types"
)

type requestHandler struct {
	middlewares []middlewares.Middleware
	offset      int
	len         int
}

// 依次执行每个路由中间件
func (h *requestHandler) Process(ctx *types.RequestCtx) (res any) {
	if h.offset >= h.len {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			ctx.Response.SetStatusCode(500)
			res = response.Error(500, fmt.Sprintln("server middleware exception:", err), nil)
			logger.Error("server middleware exception", fmt.Sprint(err))
			exception.Listener("server middleware exception", err)
		}
	}()
	han := h.middlewares[h.offset]
	h.offset++
	res = han.Process(ctx, h.Process)
	return
}
