package client

import (
	"context"
	"google.golang.org/grpc"
	"strings"
)

type clientServiceNameExtract struct {
}

type serviceName struct {
	error
	name string
}

var consumers = []string{}

func (c clientServiceNameExtract) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	if method[0] == '/' {
		method = method[1:]
	}
	name := serviceName{
		name: method[:strings.Index(method, "/")],
	}
	consumers = append(consumers, name.name)
	return name
}

func (c clientServiceNameExtract) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func GetConsumers() []string {
	return consumers
}
