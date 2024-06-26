package server

import (
	goContext "context"
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/serialize"
	"github.com/fushiliang321/go-core/helper/system"
	"github.com/hashicorp/consul/api"
	grpc1 "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"reflect"
	"strconv"
)

const HEALTHCHECK_SERVICE = "grpc.health.v1.Health"

type RegisterFunc = func(s *grpc1.Server, srv any)
type serverListen struct {
	host     string
	port     int
	server   *grpc1.Server
	listener net.Listener
}

func listen(host string, port int) *serverListen {
	address := host + ":" + strconv.Itoa(port)
	// 监听端口
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("grpc failed to listen: %v", err)
	}
	server := grpc1.NewServer(grpc1.UnaryInterceptor(func(ctx goContext.Context, req any, info *grpc1.UnaryServerInfo, handler grpc1.UnaryHandler) (resp any, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			contextDataStr := md.Get("contextData")
			if contextDataStr != nil && len(contextDataStr) > 0 && len(contextDataStr[0]) > 0 {
				var contextData map[string]any
				if serialize.JsonDecode(contextDataStr[0], &contextData) == nil {
					context.SetBatch(contextData)
				}
			}
		}
		return handler(ctx, req)
	}))

	return &serverListen{
		host:     host,
		port:     port,
		server:   server,
		listener: lis,
	}
}

func (s *serverListen) Serve() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("grpc Serve error:", fmt.Sprint(err))
			exception.Listener("grpc Serve error:", err)
		}
	}()
	serviceInfos := s.server.GetServiceInfo()
	consulConfig := consul.GetConfig()
	ip := system.GetLocalIP(consulConfig.Address)
	var serviceNames []string
	for serviceName := range serviceInfos {
		b, err := consul.RegisterServer(serviceName, "grpc", ip, s.port, &api.AgentServiceCheck{
			GRPC: fmt.Sprintf("%v:%v/%v", ip, s.port, HEALTHCHECK_SERVICE+"."+serviceName),
		})
		if b {
			serviceNames = append(serviceNames, serviceName)
			s.server.GetServiceInfo()
			event.Dispatch(event.NewRegistered(event.GrpcServerRegister, serviceName))
		} else {
			log.Printf("grpc consul register error: %v", err)
		}
	}
	if len(serviceNames) > 0 {
		s.RegisterHealthServer(serviceNames)
	}
	if err := s.server.Serve(s.listener); err != nil {
		log.Printf("grpc failed to serve: %v", err)
	}
}

// 注册服务
func (s *serverListen) RegisterServer(srv any, fun any) (res bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("grpc RegisterServer error:", fmt.Sprint(err))
			exception.Listener("grpc RegisterServer error:", err)
		}
	}()
	if fun == nil {
		logger.Warn("Not a GRPC registration function")
		return
	}
	if srv == nil {
		logger.Warn("GRPC registration parameter error")
		return
	}
	funReflect := reflect.ValueOf(fun)
	if funReflect.Kind() != reflect.Func {
		logger.Warn("Not a GRPC registration function")
		return
	}
	srvReflect := reflect.ValueOf(srv)
	if srvReflect.Kind() != reflect.Ptr {
		logger.Warn("GRPC registration parameter error")
		return
	}
	funReflect.Call([]reflect.Value{
		reflect.ValueOf(s.server),
		srvReflect,
	})
	return true
}

// 注册健康检测服务
func (s *serverListen) RegisterHealthServer(serviceNames []string) {
	if len(serviceNames) == 0 {
		return
	}
	healthserver := health.NewServer()
	for _, name := range serviceNames {
		healthserver.SetServingStatus(HEALTHCHECK_SERVICE+"."+name, healthpb.HealthCheckResponse_SERVING)
	}
	healthpb.RegisterHealthServer(s.server, healthserver)
}
