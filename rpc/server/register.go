package server

import (
	"fmt"
	"github.com/fushiliang321/go-core/consul"
	"github.com/hashicorp/consul/api"
	"log"
	"net"
	"strconv"
)

func newConsulClient() (client *api.Client, err error) {
	client, err = api.NewClient(consul.GetConfig())
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
	}
	return
}

func RegisterServer(name string, s any) {
	client, err := newConsulClient()
	if err != nil {
		log.Println("consul client error : ", err)
		return
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = name                       // 服务名称
	registration.Port, _ = strconv.Atoi(checkPort) // 服务端口
	registration.Tags = []string{}                 // tag，可以为空
	registration.Address = localIP()               // 服务 IP
	registration.Meta = map[string]string{"Protocol": "jsonrpc-http"}
	registration.Check = &api.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%s%s", registration.Address, checkPort, "/"),
		Timeout:                        "3s",
		Method:                         "POST",
		Body:                           "{\"id\":\"\",\"jsonrpc\":\"2.0\",\"method\":\"./\",\"params\":{}}",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register server error : ", err)
		return
	}
	server.Register(s)
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
