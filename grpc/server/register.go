package server

import (
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/fushiliang321/go-core/helper"
	"github.com/hashicorp/consul/api"
	grpc1 "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
	var server *grpc1.Server
	// 监听端口
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("grpc failed to listen: %v", err)
	} else {
		server = grpc1.NewServer()
		healthserver := health.NewServer()
		healthserver.SetServingStatus(HEALTHCHECK_SERVICE, healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(server, healthserver)
	}
	return &serverListen{
		host:     host,
		port:     port,
		server:   server,
		listener: lis,
	}
}

func (s *serverListen) Serve() {
	serviceInfos := s.server.GetServiceInfo()
	for serviceName := range serviceInfos {
		b, err := consul.RegisterServer(serviceName, "grpc", ip, s.port, &api.AgentServiceCheck{
			GRPC: fmt.Sprintf("%v:%v/%v", ip, s.port, HEALTHCHECK_SERVICE),
		})
		if !b {
			log.Printf("grpc consul register error: %v", err)
		}
	}
	if err := s.server.Serve(s.listener); err != nil {
		log.Printf("grpc failed to serve: %v", err)
	}
}

func (s *serverListen) RegisterServer(fun any, srv any) {
	reflect.ValueOf(fun).Call([]reflect.Value{
		reflect.ValueOf(s.server),
		reflect.ValueOf(srv),
	})
}
