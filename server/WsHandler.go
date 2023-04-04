package server

import "github.com/valyala/fasthttp"

type WsController struct{}

func (ws *WsController) Handler(ctx *fasthttp.RequestCtx) any {
	return ""
}
