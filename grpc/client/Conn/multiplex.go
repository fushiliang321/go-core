package client

import (
	goContext "context"
	"errors"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/serialize"
	context2 "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"sync"
	"sync/atomic"
)

type Multiplex struct {
	*Base
	multiplexNum  atomic.Uint32 //连接已复用次数
	serviceNode   *consul.ServiceNode
	isCloseClient bool //客户端是否关闭
	sync.RWMutex
}

func NewMultiplex(serviceName string) *Multiplex {
	return &Multiplex{
		Base: &Base{
			serviceName: serviceName,
		},
	}
}

func (c *Multiplex) getConn() (conn *grpc.ClientConn, err error) {
	c.RLock()
	if c.isCloseClient {
		c.RUnlock()
		return nil, errors.New("The client is closed ")
	}

	if c.cc != nil && // 连接存在
		c.serviceNode != nil && !c.serviceNode.IsRemove && // 服务节点存在且没有被移除
		(maxMultiplexNum == 0 || c.multiplexNum.Add(1) < maxMultiplexNum) { // 连接未达到最大复用次数
		conn = c.cc
		c.RUnlock()
		return
	}
	c.RUnlock()
	c.Lock()
	defer c.Unlock()
	if c.cc != nil {
		oldConn := c.cc
		go func() {
			oldConn.WaitForStateChange(context2.Background(), connectivity.Connecting)
			oldConn.Close()
		}()
	}

	conn, c.serviceNode, err = dial(c.serviceName)
	if err != nil {
		c.cc = nil
		c.serviceNode = nil
		return nil, err
	}
	c.cc = conn
	c.multiplexNum.Store(0)
	return conn, nil
}

func (c *Multiplex) Invoke(ctx goContext.Context, method string, args, reply any, opts ...grpc.CallOption) (err error) {
	var con *grpc.ClientConn
	con, err = c.getConn()
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = goContext.Background()
	}
	contextData := context.GetAll()
	if len(contextData) > 0 {
		var str string
		if str, err = serialize.JsonEncode(contextData); err == nil {
			ctx = metadata.AppendToOutgoingContext(ctx, "contextData", str)
		}
	}
	err = con.Invoke(ctx, method, args, reply, opts...)
	if err != nil {
		logger.Warn("grpc client Invoke error：", err)
	}
	return err
}
