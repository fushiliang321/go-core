package client

import (
	"context"
	"fmt"
	"github.com/fushiliang321/go-core/exception"
	client "github.com/fushiliang321/go-core/grpc/client/Conn"
	"github.com/fushiliang321/go-core/helper/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"reflect"
	"sync"
)

type GenerateFun[t any] func(isMultiplex ...bool) t

var (
	multiplexConns = map[any]client.Interface{}
	ctx            = context.Background()
)

func NewClient[t any](serviceName string, fun func(cc grpc.ClientConnInterface) t) GenerateFun[t] {
	return func(isMultiplex ...bool) t {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("grpc call exception:", fmt.Sprint(err))
				exception.Listener("grpc call exception:", err)
			}
		}()
		var multiplex bool
		if len(isMultiplex) > 0 {
			multiplex = isMultiplex[0]
		}
		conn, err := GetConn(serviceName, multiplex)
		if err != nil {
			logger.Warn("grpc newClient Error:["+serviceName+"]", err)
			exception.Listener("grpc newClient Error:["+serviceName+"]", err)
			return fun(ErrCon{
				error: err,
			})
		}
		_client := fun(conn)
		if multiplex {
			multiplexConns[_client] = conn
		}
		return _client
	}
}

func Auto[t any](fun func(cc grpc.ClientConnInterface) t) GenerateFun[t] {
	_client := fun(clientServiceNameExtract{})
	name := ""
	if reflect.ValueOf(_client).NumMethod() > 0 {
		reflectMethod := reflect.ValueOf(_client).Method(0)
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
		delete(multiplexConns, client)
		conn.Close()
	}
}

// 关闭复用客户端，等待所有请求结束后再关闭
func CloseAwait(client any) {
	wg := &sync.WaitGroup{}
	if conn, ok := multiplexConns[client]; ok {
		delete(multiplexConns, client)
		go func() {
			wg.Add(1)
			defer wg.Done()
			conn.WaitForStateChange(ctx, connectivity.Connecting)
			conn.Close()
		}()
	}
	wg.Wait()
}
