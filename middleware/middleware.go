package middleware

import (
	"github.com/fushiliang321/go-core/config/middlewares"
	"github.com/fushiliang321/go-core/router/types"
)

func GetCoreMiddlewares(handler types.RequestHandler) *middlewares.Middlewares {
	return &middlewares.Middlewares{
		Http: []middlewares.Middleware{
			&HttpCoreMiddleware{
				handler,
			},
		},
		WS: []middlewares.Middleware{
			&WebsocketCoreMiddleware{},
		},
	}
}
