package client

import (
	context2 "context"
	"fmt"
	"github.com/fushiliang321/go-core/config/jsonRpcHttp"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/event/handles/core"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	string2 "github.com/fushiliang321/go-core/helper/string"
	"github.com/fushiliang321/jsonrpc"
)

type Client struct {
	serverName      string
	serverNameSnake string
}

var DefaultContext func() context2.Context

func init() {
	go func() {
		core.AwaitStartFinish()
		if jsonRpcHttp.Get().DefaultContext != nil {
			DefaultContext = jsonRpcHttp.Get().DefaultContext
		}
	}()
}

func New(server string) *Client {
	return &Client{
		serverName:      server,
		serverNameSnake: string2.SnakeString(server),
	}
}

func (c *Client) Call(method string, params any, res any) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("rpc call error：", fmt.Sprint(err))
			exception.Listener("rpc call", err)
		}
	}()
	rpcClient, err := newRpcClient(c.serverName)
	if err == nil {
		err = rpcClient.Call(c.serverNameSnake+"/"+method, params, res, false, context.GetAll())
		if err != nil {
			logger.Warn("rpc rpcClient error:", err.Error())
		}
	}
	return err
}

func (c *Client) CallWithContext(ctx context2.Context, method string, params any, res any) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("rpc call error：", fmt.Sprint(err))
			exception.Listener("rpc call", err)
		}
	}()
	rpcClient, err := newRpcClient(c.serverName)
	if err == nil {
		if ctx == nil && DefaultContext != nil {
			ctx = DefaultContext()
		}
		err = rpcClient.CallWithContext(ctx, c.serverNameSnake+"/"+method, params, res, false, context.GetAll())
		if err != nil {
			logger.Warn("rpc rpcClient error:", err.Error())
		}
	}
	return err
}

func newRpcClient(name string) (jsonrpc.ClientInterface, error) {
	node, err := consul.GetNode(name+"Service", consul.HttpProtocol)
	if err != nil {
		logger.Warn("rpc newClient Error: ["+name+"]", err.Error())
		exception.Listener("rpc newClient Error: ["+name+"]", err)
		return nil, err
	}
	return jsonrpc.NewClient(node.Protocol, node.Address, node.Port)
}

func CallWithContext(ctx context2.Context, server string, method string, params any, res any) error {
	return New(server).CallWithContext(ctx, method, params, res)
}

func Call(server string, method string, params any, res any) error {
	return New(server).Call(method, params, res)
}
