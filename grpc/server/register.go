package server

import (
	goContext "context"
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/context"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
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

var ip string

func init() {
	ip = helper.GetLocalIP()
}

func listen(host string, port int) *serverListen {
	address := host + ":" + strconv.Itoa(port)
	// 监听端口
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("grpc failed to listen: %v", err)
	}
	server := grpc1.NewServer(grpc1.UnaryInterceptor(func(ctx goContext.Context, req interface{}, info *grpc1.UnaryServerInfo, handler grpc1.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			contextDataStr := md.Get("contextData")
			if contextDataStr != nil && len(contextDataStr) > 0 && len(contextDataStr[0]) > 0 {
				var contextData map[string]any
				err := helper.JsonDecode(contextDataStr[0], &contextData)
				if err == nil {
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
			fmt.Println("grpc Serve error:", err)
			exception.Listener("grpc Serve error:", err)
		}
	}()
	serviceInfos := s.server.GetServiceInfo()
	for serviceName := range serviceInfos {
		b, err := consul.RegisterServer(serviceName, "grpc", ip, s.port, &api.AgentServiceCheck{
			GRPC: fmt.Sprintf("%v:%v/%v", ip, s.port, HEALTHCHECK_SERVICE),
		})
		if !b {
			log.Printf("grpc consul register error: %v", err)
		}
	}
	s.RegisterHealthServer()
	if err := s.server.Serve(s.listener); err != nil {
		log.Printf("grpc failed to serve: %v", err)
	}
}

func (s *serverListen) RegisterServer(srv any, fun any) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("grpc RegisterServer error:", err)
			exception.Listener("grpc RegisterServer error:", err)
		}
	}()
	var reflectValue reflect.Value
	if reflect.TypeOf(fun).Kind() == reflect.Ptr {
		reflectValue = reflect.ValueOf(fun).Elem()
	} else {
		reflectValue = reflect.ValueOf(fun)
	}
	if reflectValue.Kind() != reflect.Func {
		fmt.Println("Not a GRPC registration function")
		return
	}
	if reflect.TypeOf(srv).Kind() != reflect.Ptr {
		fmt.Println("GRPC registration parameter error")
		return
	}
	reflectValue.Call([]reflect.Value{
		reflect.ValueOf(s.server),
		reflect.ValueOf(srv),
	})
}

func (s *serverListen) RegisterHealthServer() {
	healthserver := health.NewServer()
	healthserver.SetServingStatus(HEALTHCHECK_SERVICE, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s.server, healthserver)
}
