package middlewares

import (
	"gitee.com/zvc/go-core/router/types"
	"github.com/valyala/fasthttp"
)

type Middlewares struct {
	Http []Middleware
	WS   []Middleware
}

type Middleware interface {
	Process(ctx *fasthttp.RequestCtx, handler types.RequestHandler) any
}

var middlewares = &Middlewares{
	Http: []Middleware{},
	WS:   []Middleware{},
}

func Set(config *Middlewares) {
	middlewares = config
}

func Get() *Middlewares {
	return middlewares
}
