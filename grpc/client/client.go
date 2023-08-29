package client

import (
	"context"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"google.golang.org/grpc"
	"reflect"
)

type ClientGenerateFun[t any] func(isMultiplex ...bool) t

var (
	multiplexConns = map[any]*ClientConn{}
	ctx            = context.Background()
)

func NewClient[t any](serviceName string, fun func(cc grpc.ClientConnInterface) t) ClientGenerateFun[t] {
	return func(isMultiplex ...bool) t {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("grpc call exception:", err)
				exception.Listener("grpc call exception:", err)
			}
		}()
		var multiplex bool
		if len(isMultiplex) > 0 {
			multiplex = isMultiplex[0]
		}
		conn, err := GetConn(serviceName, multiplex)
		var client t
		if err != nil {
			logger.Warn("grpc newClient Error:["+serviceName+"]", err)
			exception.Listener("grpc newClient Error:["+serviceName+"]", err)
			return fun(ErrCon{
				error: err,
			})
		}

		client = fun(conn)
		if multiplex {
			multiplexConns[client] = conn
		}
		return client
	}
}

func ClientAuto[t any](fun func(cc grpc.ClientConnInterface) t) ClientGenerateFun[t] {
	client := fun(clientServiceNameExtract{})
	name := ""
	if reflect.ValueOf(client).NumMethod() > 0 {
		reflectMethod := reflect.ValueOf(client).Method(0)
		res := reflectMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.New(reflectMethod.Type().In(1).Elem()),
		})
		err := res[1].Interface()
		if serviceNameType, ok := err.(serviceName); ok {
			name = serviceNameType.name
		}
	}
	return NewClient[t](name, fun)
}

// 关闭复用客户端
func Close(client any) {
	if conn, ok := multiplexConns[client]; ok {
		defer delete(multiplexConns, client)
		conn.isCloseClient = true
		conn.Close()
	}
}

// 关闭复用客户端，等待所有请求结束后再关闭
func CloseAwait(client any) {
	if conn, ok := multiplexConns[client]; ok {
		defer delete(multiplexConns, client)
		conn.Lock()
		defer conn.Unlock()
		conn.isCloseClient = true
		conn.Close()
	}
}
