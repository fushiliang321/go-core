package client

import (
	"context"
	"google.golang.org/grpc"
	"strings"
)

type clientServiceNameExtract struct {
}

var consumers = []string{}

func (c clientServiceNameExtract) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	if method[0] == '/' {
		method = method[1:]
	}
	serviceNamePtr := ctx.Value("serviceName").(*string)
	*serviceNamePtr = method[:strings.Index(method, "/")]
	consumers = append(consumers, *serviceNamePtr)
	ctx = context.WithValue(ctx, "key1", "modify from v21")
	return nil
}
func (c clientServiceNameExtract) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func GetConsumers() []string {
	return consumers
}