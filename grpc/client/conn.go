package client

import (
	"fmt"
	"github.com/fushiliang321/go-core/exception"
	client "github.com/fushiliang321/go-core/grpc/client/Conn"
	"github.com/fushiliang321/go-core/helper/logger"
)

func GetConn(serviceName string, multiplex bool) (con client.Interface, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("server middleware exception", fmt.Sprint(err))
			exception.Listener("grpc conn exception", err)
		}
	}()

	if multiplex {
		con = client.NewMultiplex(serviceName)
	} else {
		con = client.New(serviceName)
	}
	return
}
