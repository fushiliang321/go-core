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
	"sync"
)

type ClientConn struct {
	serviceName      string
	cc               *grpc.ClientConn
	ccInUse          bool      //连接是否被使用
	multiplex        bool      //是否复用连接
	currentLimitChan chan byte //限流通道
	sync.RWMutex
}

func dial(serviceName string) (*grpc.ClientConn, error) {
	defer func() {
		exception.Listener("grpc dial exception", recover())
	}()
	node, err := consul.GetNode(serviceName, consul.GrpcProtocol)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(node.Address+":"+node.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (cc *ClientConn) Invoke(ctx goContext.Context, method string, args, reply interface{}, opts ...grpc.CallOption) (err error) {
	var con *grpc.ClientConn
	if !cc.multiplex {
		//不复用连接的情况下 每次调用都会重新连接
		//第一次调用不需要重新连接，直接使用连接，避免浪费
		//复用连接的情况下不需要自动重新连接，避免重连后会忘记去手动关闭连接
		if cc.ccInUse {
			con, err = dial(cc.serviceName)
			if err != nil {
				return err
			}
		} else {
			cc.ccInUse = true
			con = cc.cc
		}
		defer con.Close()
	} else {
		cc.currentLimitChan <- <-cc.currentLimitChan
		cc.RLock()
		defer cc.RUnlock()
		con = cc.cc
	}
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
	err = con.Invoke(ctx, method, args, reply, opts...)
	if err != nil {
		fmt.Println("grpc client Invoke error：", err)
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
	node, err := consul.GetNode(serviceName, consul.GrpcProtocol)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(node.Address+":"+node.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	currentLimitChan := make(chan byte, 1)
	if multiplex {
		currentLimitChan <- 0
	}
	return &ClientConn{
		serviceName:      serviceName,
		cc:               conn,
		multiplex:        multiplex,
		currentLimitChan: currentLimitChan,
	}, nil
}
