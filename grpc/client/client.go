package client

import (
	"github.com/fushiliang321/go-core/exception"
	"google.golang.org/grpc"
	"reflect"
	"strings"
)

type ClientGenerateFun[t any] func(isMultiplex ...bool) t

var multiplexConns = map[any]*ClientConn{}

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
		if err != nil {
			exception.Listener("grpc newClient Error: ["+serviceName+"]", err)
			return *new(t)
		}
		client := fun(conn)
		if multiplex {
			multiplexConns[client] = conn
		}
		return client
	}
}

func ClientAuto[t any](fun func(cc grpc.ClientConnInterface) t) ClientGenerateFun[t] {
	var t1 = new(t)
	typeStr := reflect.TypeOf(t1).String()
	lastIndex := strings.LastIndexAny(reflect.TypeOf(t1).String(), ".")
	serviceName := typeStr[lastIndex+1 : len(typeStr)-6]
	return NewClient[t](serviceName, fun)
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
