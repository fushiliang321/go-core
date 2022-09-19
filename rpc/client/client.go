package client

import (
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	jsonrpc "github.com/iloveswift/go-jsonrpc"
	"log"
)

type Client struct {
	serverName      string
	serverNameSnake string
}

func New(server string) *Client {
	return &Client{
		serverName:      server,
		serverNameSnake: helper.SnakeString(server),
	}
}

func (c *Client) Call(method string, params any, res any) (err error) {
	defer func() {
		exception.Listener("rpc call", recover())
	}()
	rpcClient, err := newRpcClient(c.serverName)
	if err == nil {
		err = rpcClient.Call(c.serverNameSnake+"/"+method, params, res, false)
		if err != nil {
			log.Println("rpcClient error", err)
		}
	}
	return
}

func newRpcClient(name string) (jsonrpc.ClientInterface, error) {
	node, err := consul.GetNode(name + "Service")
	if err != nil {
		exception.Listener("newClient Error: ["+name+"]", err)
		return nil, err
	}
	return jsonrpc.NewClient(node.Protocol, node.Address, node.Port)
}

func Call(server string, method string, params any, res any) (err error) {
	return New(server).Call(method, params, res)
}
