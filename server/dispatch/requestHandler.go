package dispatch

import (
	"core/config/middlewares"
	"core/exception"
	"core/helper"
	"fmt"
	"github.com/valyala/fasthttp"
)

type requestHandler struct {
	middlewares []middlewares.Middleware
	offset      int
	len         int
}

func (h *requestHandler) Process(ctx *fasthttp.RequestCtx) (res any) {
	if h.offset >= h.len {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			ctx.Response.SetStatusCode(500)
			res = helper.Error(500, fmt.Sprintln("server middleware exception:", err), nil)
			exception.Listener("server middleware exception", err)
		}
	}()
	han := h.middlewares[h.offset]
	h.offset++
	res = han.Process(ctx, h.Process)
	return
}
