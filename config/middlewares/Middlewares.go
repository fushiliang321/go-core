package middlewares

import (
	"github.com/fushiliang321/go-core/router/types"
)

type (
	Middlewares struct {
		Http []Middleware
		WS   []Middleware
	}
	Middleware interface {
		Process(ctx *types.RequestCtx, handler types.RequestHandler) any
	}
)

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
