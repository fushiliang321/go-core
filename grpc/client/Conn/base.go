package client

import (
	goContext "context"
	"errors"
	"fmt"
	grpcConfig "github.com/fushiliang321/go-core/config/grpc"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/event/handles/core"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var maxMultiplexNum uint32 //连接最大复用次数

type Base struct {
	serviceName string
	cc          *grpc.ClientConn
}

func init() {
	core.AwaitStartFinish()
	config := grpcConfig.Get()
	maxMultiplexNum = config.ConnectMaxMultiplexNum
}

func (c *Base) NewStream(ctx goContext.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if c.cc == nil {
		return nil, errors.New("连接不存在")
	}
	return c.cc.NewStream(ctx, desc, method, opts...)
}

func (c *Base) GetState() connectivity.State {
	if c.cc == nil {
		return connectivity.Shutdown
	}
	return c.cc.GetState()
}

func (c *Base) Target() string {
	if c.cc == nil {
		return ""
	}
	return c.cc.Target()
}

func (c *Base) WaitForStateChange(ctx goContext.Context, sourceState connectivity.State) bool {
	if c.cc == nil {
		return false
	}
	return c.cc.WaitForStateChange(ctx, sourceState)
}

func (c *Base) Connect() {
	if c.cc == nil {
		return
	}
	c.cc.Connect()
}

func (c *Base) Close() error {
	defer func() {
		c.cc = nil
		if err := recover(); err != nil {
			logger.Error("conn close error:", fmt.Sprint(err))
		}
	}()
	if c.cc == nil {
		return nil
	}
	return c.cc.Close()
}

func (c *Base) ResetConnectBackoff() {
	if c.cc == nil {
		return
	}
	c.cc.ResetConnectBackoff()
}

func (c *Base) GetMethodConfig(method string) grpc.MethodConfig {
	if c.cc == nil {
		return grpc.MethodConfig{}
	}
	return c.cc.GetMethodConfig(method)
}

func dial(serviceName string) (conn *grpc.ClientConn, node *consul.ServiceNode, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("grpc dial exception")
			logger.Error("grpc dial exception:", fmt.Sprint(e))
			exception.Listener("grpc dial exception", e)
		}
	}()
	node, err = consul.GetNode(serviceName, consul.GrpcProtocol)
	if err != nil {
		return nil, nil, err
	}
	conn, err = grpc.NewClient(node.Address+":"+node.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, node, err
}
