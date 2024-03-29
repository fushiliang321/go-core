package client

import (
	"context"
	"google.golang.org/grpc"
)

type ErrCon struct {
	error
	connType
}

func (c ErrCon) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return c
}
