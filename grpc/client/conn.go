package client

import (
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/exception"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetConn(serviceName string) (*grpc.ClientConn, error) {
	defer func() {
		exception.Listener("grpc conn exception", recover())
	}()
	node, err := consul.GetNode(serviceName, "grpc")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(node.Address+":"+node.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, nil
}
