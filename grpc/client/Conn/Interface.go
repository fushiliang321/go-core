package client

import (
	goContext "context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Interface interface {
	grpc.ClientConnInterface
	GetState() connectivity.State
	WaitForStateChange(ctx goContext.Context, sourceState connectivity.State) bool
	Close() error
}
