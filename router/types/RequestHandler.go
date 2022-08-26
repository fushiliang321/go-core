package types

import "github.com/valyala/fasthttp"

type RequestHandler = func(ctx *fasthttp.RequestCtx) any
