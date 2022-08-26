package middleware

import (
	"core/config/middlewares"
	"core/router/types"
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
