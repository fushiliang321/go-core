package client

import (
	"core/consul"
	"core/exception"
	"core/helper"
	jsonrpc "github.com/iloveswift/go-jsonrpc"
	"log"
)

func newRpcClient(name string) (jsonrpc.ClientInterface, error) {
	node, err := consul.GetNode(name + "Service")
	if err != nil {
		exception.Listener("newClient Error: ["+name+"]", err)
		return nil, err
	}
	return jsonrpc.NewClient(node.Protocol, node.Address, node.Port)
}

func Call(server string, method string, params any, res any) (err error) {
	defer func() {
		exception.Listener("rpc call", recover())
	}()
	rpcClient, err := newRpcClient(server)
	if err == nil {
		err = rpcClient.Call(helper.SnakeString(server)+"/"+method, params, res, false)
		if err != nil {
			log.Println("rpcClient error", err)
		}
	}
	return
}
