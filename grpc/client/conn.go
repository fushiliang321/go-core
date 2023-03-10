package client

import (
	goContext "context"
	"errors"
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
	multiplexNum     uint      //连接已复用次数
	currentLimitChan chan byte //限流通道
	sync.RWMutex
}

const maxMultiplexNum = 1000 //连接最大复用次数

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

	if cc.multiplex {
		//连接复用
		cc.currentLimitChan <- <-cc.currentLimitChan
		cc.RLock()
		defer cc.RUnlock()

		if cc.cc == nil {
			cc.cc, err = dial(cc.serviceName)
			if err != nil {
				return err
			}
			cc.multiplexNum = 0
		}

		con = cc.cc

		defer func() {
			cc.multiplexNum++
			if cc.multiplexNum < maxMultiplexNum {
				return
			}
			//超出复用次数就关闭连接
			cc.cc.Close()
			cc.cc = nil
		}()
	} else {
		//不复用连接的情况下 每次调用都会重新连接
		cc.cc, err = dial(cc.serviceName)
		if err != nil {
			return err
		}
		con = cc.cc
		defer func() {
			cc.cc = nil
			con.Close()
		}()
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
	if cc.cc == nil {
		return nil, errors.New("连接不存在")
	}
	return cc.cc.NewStream(ctx, desc, method, opts...)
}
func (cc *ClientConn) GetState() connectivity.State {
	if cc.cc == nil {
		return connectivity.Shutdown
	}
	return cc.cc.GetState()
}
func (cc *ClientConn) Target() string {
	if cc.cc == nil {
		return ""
	}
	return cc.cc.Target()
}
func (cc *ClientConn) WaitForStateChange(ctx goContext.Context, sourceState connectivity.State) bool {
	if cc.cc == nil {
		return false
	}
	return cc.cc.WaitForStateChange(ctx, sourceState)
}
func (cc *ClientConn) Connect() {
	if cc.cc == nil {
		return
	}
	cc.cc.Connect()
}
func (cc *ClientConn) Close() error {
	if cc.cc == nil {
		return nil
	}
	err := cc.cc.Close()
	if err != nil {
		return err
	}
	cc.cc = nil
	return nil
}
func (cc *ClientConn) ResetConnectBackoff() {
	if cc.cc == nil {
		return
	}
	cc.cc.ResetConnectBackoff()
}
func (cc *ClientConn) GetMethodConfig(method string) grpc.MethodConfig {
	if cc.cc == nil {
		return grpc.MethodConfig{}
	}
	return cc.cc.GetMethodConfig(method)
}

func GetConn(serviceName string, multiplex bool) (*ClientConn, error) {
	defer func() {
		exception.Listener("grpc conn exception", recover())
	}()
	var currentLimitChan chan byte
	if multiplex {
		currentLimitChan = make(chan byte, 1)
		currentLimitChan <- 0
	}
	return &ClientConn{
		serviceName:      serviceName,
		cc:               nil,
		multiplex:        multiplex,
		currentLimitChan: currentLimitChan,
	}, nil
}
