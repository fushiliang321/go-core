package client

import (
	"github.com/fushiliang321/go-core/exception"
	"google.golang.org/grpc"
	"reflect"
	"strings"
)

func NewClient[t any](serviceName string, fun func(cc grpc.ClientConnInterface) t) func() t {
	return func() t {
		defer func() {
			exception.Listener("grpc call exception", recover())
		}()
		conn, err := GetConn(serviceName)
		if err != nil {
			exception.Listener("grpc newClient Error: ["+serviceName+"]", err)
			return *new(t)
		}
		return fun(conn)
	}
}

func ClientAuto[t any](fun func(cc grpc.ClientConnInterface) t) func() t {
	var t1 = new(t)
	typeStr := reflect.TypeOf(t1).String()
	lastIndex := strings.LastIndexAny(reflect.TypeOf(t1).String(), ".")
	serviceName := typeStr[lastIndex+1 : len(typeStr)-6]
	return NewClient[t](serviceName, fun)
}
