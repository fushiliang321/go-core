package request

import (
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/jsonRpcHttp/context"
	"github.com/fushiliang321/go-core/router/types"
	"github.com/fushiliang321/go-core/rpc"
	"github.com/savsgio/gotils/strconv"
	"net"
	"reflect"
)

var ipHeaderKeys = []string{ //请求头中可以获取到客户端ip的字段
	"client-ip",
	"real_client_ip",
	"x-forwarded-for",
	"x-real-ip",
	"x-true-ip",
	"wl-proxy-client-ip",
	"x-client-ip",
}

// 获取客户端ip
func ClientIP(ctx *types.RequestCtx) (string, uint8) {
	var (
		ip                net.IP
		cip               []byte
		isPublicNetworkIp = false //是否为公网ip
	)

	for _, key := range ipHeaderKeys {
		cip = ctx.Request.Header.Peek(key)
		if cip != nil {
			ip = net.ParseIP(strconv.B2S(cip))
			if ip != nil && !IsIntranetIp(ip) {
				isPublicNetworkIp = true
				break
			}
		}
	}

	if !isPublicNetworkIp {
		ip = ctx.RemoteIP()
	}

	if ip == nil {
		return "", 0
	}
	s := ip.String()
	return s, helper.IpType(s)
}

// 判断是否为内网ip
func IsIntranetIp(ip net.IP) (b bool) {
	ip4 := ip.To4()
	if ip4 == nil {
		return
	}
	switch ip4[0] {
	case 10:
		b = true
	case 172:
		b = ip4[1] >= 16 && ip4[1] <= 31
	case 169:
		b = ip4[1] == 254
	case 192:
		b = ip4[1] == 168
	}
	return
}

// 获取客户端地址（ip+port）
func ClientAddr(ctx *types.RequestCtx) string {
	ip, v := ClientIP(ctx)
	if ip == "" {
		return ""
	}
	port := ctx.Request.Header.Peek("client-port")
	if port == nil {
		port = ctx.Request.Header.Peek("real_client_port")
	}
	if port == nil {
		if ip == ctx.RemoteIP().String() {
			return ctx.RemoteAddr().String()
		}
		return ip
	}
	if v == 4 {
		return ip + ":" + strconv.B2S(port)
	}
	return "[" + ip + "]:" + strconv.B2S(port)
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
		err := helper.MapToStruc[string](mapData, &rpcRequestData)
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

func Input[T any](ctx *types.RequestCtx, key string, defaultVal ...T) (*T, error) {
	var v T
	_reflect := reflect.New(reflect.TypeOf(v)).Elem()
	if len(defaultVal) > 0 {
		_reflect.Set(reflect.ValueOf(defaultVal[0]))
	}
	v = _reflect.Interface().(T)
	if err := ctx.InputAssign(key, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
