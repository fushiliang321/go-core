package middleware

import (
	"gitee.com/zvc/go-core/config/middlewares"
	"gitee.com/zvc/go-core/router/types"
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
