package client

import (
	goContext "context"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/serialize"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Conn struct {
	*Base
}

func New(serviceName string) *Conn {
	return &Conn{
		Base: &Base{
			serviceName: serviceName,
		},
	}
}

func (c *Conn) Invoke(ctx goContext.Context, method string, args, reply any, opts ...grpc.CallOption) (err error) {
	var con *grpc.ClientConn
	con, _, err = dial(c.serviceName)
	if err != nil {
		return err
	}
	c.cc = con
	defer func() {
		c.cc = nil
		con.Close()
	}()
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
