package client

import (
	"context"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/exception"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConn struct {
	cc *grpc.ClientConn
}

func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	err := cc.cc.Invoke(ctx, method, args, reply, opts...)
	cc.Close()
	return err
}
func (cc *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return cc.cc.NewStream(ctx, desc, method, opts...)
}
func (cc *ClientConn) GetState() connectivity.State {
	return cc.cc.GetState()
}
func (cc *ClientConn) Target() string {
	return cc.cc.Target()
}
func (cc *ClientConn) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {
	return cc.cc.WaitForStateChange(ctx, sourceState)
}
func (cc *ClientConn) Connect() {
	cc.cc.Connect()
}
func (cc *ClientConn) Close() error {
	return cc.cc.Close()
}
func (cc *ClientConn) ResetConnectBackoff() {
	cc.cc.ResetConnectBackoff()
}
func (cc *ClientConn) GetMethodConfig(method string) grpc.MethodConfig {
	return cc.cc.GetMethodConfig(method)
}

func GetConn(serviceName string) (grpc.ClientConnInterface, error) {
	defer func() {
		exception.Listener("grpc conn exception", recover())
	}()
	node, err := consul.GetNode(serviceName, "grpc")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(node.Address+":"+node.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ClientConn{
		cc: conn,
	}, nil
}
