package response

import (
	"github.com/fushiliang321/go-core/helper/system"
	"github.com/fushiliang321/go-core/helper/types"
	types2 "github.com/fushiliang321/go-core/router/types"
)

// 响应成功数据
func Success(msg string, datas ...any) (res *types.Result) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	return &types.Result{
		Code:    1,
		Msg:     msg,
		Data:    data,
		ErrCode: 0,
		Service: system.AppName(),
	}
}

// 响应错误数据
func Error(errCode int, msg string, datas ...any) (res *types.Result) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	return &types.Result{
		Code:    0,
		Msg:     msg,
		Data:    data,
		ErrCode: errCode,
		Service: system.AppName(),
	}
}

// 响应错误数据
func ErrorResponse(ctx *types2.RequestCtx, errCode int, msg string, datas ...any) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	ctx.Response.SetStatusCode(errCode)
	res := types.Result{}
	res.Error(errCode, msg, data)
	_, _ = ctx.Write(res.JsonMarshal())
}
