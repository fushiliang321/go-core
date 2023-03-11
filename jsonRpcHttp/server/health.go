package server

import (
	"github.com/fushiliang321/jsonrpc/common"
	"github.com/valyala/fasthttp"
)

type Health struct{}

var (
	resultSuccess = "success"
	resultError   = "error"
	internalErr   = common.NewInternalErr(resultError, ErrorResponse{
		Code: fasthttp.StatusGone,
		Text: resultError,
	})
)

func (s *Health) Check(params *registerInfo) (*string, error) {
	if len(serviceRegistrations) == 0 {
		//没有注册的服务
		return &resultError, internalErr
	}
	registration, ok := serviceRegistrations[params.Name]
	if !ok {
		//检测的服务没有注册
		return &resultError, internalErr
	}
	if registration.Protocol != params.Protocol {
		//服务的协议不一致
		return &resultError, internalErr
	}
	return &resultSuccess, nil
}
