package helper

import (
	"github.com/fushiliang321/go-core/jsonRpcHttp/context"
	"github.com/fushiliang321/go-core/rpc"
	"github.com/valyala/fasthttp"
	"net"
)

// 获取客户端ip
func ClientIP(ctx *fasthttp.RequestCtx) (string, uint8) {
	cip := ctx.Request.Header.Peek("client-ip")
	var ip net.IP
	if cip != nil {
		ip = net.ParseIP(string(cip))
	} else {
		ip = ctx.RemoteIP()
	}
	if ip == nil {
		return "", 0
	}
	s := ip.String()
	return s, IpType(s)
}

// 获取客户端地址（ip+port）
func ClientAddr(ctx *fasthttp.RequestCtx) string {
	ip, v := ClientIP(ctx)
	if ip == "" {
		return ""
	}
	port := ctx.Request.Header.Peek("client-port")
	if port == nil {
		if ip == ctx.RemoteIP().String() {
			return ctx.RemoteAddr().String()
		}
		return ip
	}
	if v == 4 {
		return ip + ":" + string(port)
	}
	return "[" + ip + "]:" + string(port)
}

const internalRequestKey = "internalRequest"

// 获取rpc上下文请求数据
func RpcRequestData() (rpcRequestData rpc.RpcRequestData) {
	data := context.Get(internalRequestKey)
	if data == nil {
		return
	}
	switch data.(type) {
	case rpc.RpcRequestData:
		rpcRequestData = data.(rpc.RpcRequestData)
	case map[string]any:
		mapData, ok := data.(map[string]any)
		if !ok {
			return
		}
		err := MapToStruc[string](mapData, &rpcRequestData)
		if err != nil {
			return
		}
	default:
		return
	}
	return rpcRequestData
}

// 设置rpc上下文请求数据
func SetRpcRequestData(rpcRequestData rpc.RpcRequestData) {
	context.Set(internalRequestKey, rpcRequestData)
}
