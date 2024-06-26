package client

import (
	"context"
	"google.golang.org/grpc"
	"strings"
)

type (
	clientServiceNameExtract struct {
		connType
	}
	serviceName struct {
		error
		name string
	}
)

var consumers = []string{}

func GetConsumers() []string {
	return consumers
}

func (c clientServiceNameExtract) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	if method[0] == '/' {
		method = method[1:]
	}
	name := serviceName{
		name: method[:strings.Index(method, "/")],
	}
	consumers = append(consumers, name.name)
	return name
}
