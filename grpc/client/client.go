package client

import (
	"context"
	"github.com/fushiliang321/go-core/exception"
	"google.golang.org/grpc"
	"reflect"
)

type ClientGenerateFun[t any] func(isMultiplex ...bool) t

var multiplexConns = map[any]*ClientConn{}

var ctx context.Context

func init() {
	ctx = context.Background()
}

func NewClient[t any](serviceName string, fun func(cc grpc.ClientConnInterface) t) ClientGenerateFun[t] {
	return func(isMultiplex ...bool) t {
		defer func() {
			exception.Listener("grpc call exception:", recover())
		}()
		var multiplex bool
		if len(isMultiplex) > 0 {
			multiplex = isMultiplex[0]
		}
		conn, err := GetConn(serviceName, multiplex)
		var client t
		if err != nil {
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
		conn.Close()
	}
}

// 关闭复用客户端，等待所有请求结束后再关闭
func CloseAwait(client any) {
	if conn, ok := multiplexConns[client]; ok {
		defer delete(multiplexConns, client)
		v := <-conn.currentLimitChan
		conn.Lock()
		defer conn.Unlock()
		conn.currentLimitChan <- v
		conn.Close()
	}
}
