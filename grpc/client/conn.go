package client

import (
	goContext "context"
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ClientConn struct {
	cc        *grpc.ClientConn
	multiplex bool //是否复用连接
}

func (cc *ClientConn) Invoke(ctx goContext.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	if ctx == nil {
		ctx = goContext.Background()
	}
	contextData := context.GetAll()
	if contextData != nil && len(contextData) > 0 {
		str, err := helper.JsonEncode(contextData)
		if err == nil {
			ctx = metadata.AppendToOutgoingContext(ctx, "contextData", str)
		}
	}
	err := cc.cc.Invoke(ctx, method, args, reply, opts...)
	if err != nil {
		fmt.Println("grpc client Invoke error：", err)
	}
	if !cc.multiplex {
		cc.Close()
	}
	return err
}
func (cc *ClientConn) NewStream(ctx goContext.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return cc.cc.NewStream(ctx, desc, method, opts...)
}
func (cc *ClientConn) GetState() connectivity.State {
	return cc.cc.GetState()
}
func (cc *ClientConn) Target() string {
	return cc.cc.Target()
}
func (cc *ClientConn) WaitForStateChange(ctx goContext.Context, sourceState connectivity.State) bool {
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

func GetConn(serviceName string, multiplex bool) (*ClientConn, error) {
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
		cc:        conn,
		multiplex: multiplex,
	}, nil
}
