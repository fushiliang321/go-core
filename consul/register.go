package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

func newConsulClient() (client *api.Client, err error) {
	client, err = api.NewClient(GetConfig())
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
	}
	return
}

func RegisterServer(name string, protocol string, address string, port int, check *api.AgentServiceCheck) (b bool, err error) {
	client, err := newConsulClient()
	if err != nil {
		log.Println("consul client error : ", err)
		return
	}
	if check != nil {
		check = setServiceCheckDefaultValue(check)
	}

	registration := &api.AgentServiceRegistration{
		Name: name, // 服务名称
		Port: port, // 服务端口
		ID:   name, //服务id
		// todo 需要完善服务唯一id生成规则
		Tags:    []string{protocol}, // tag，可以为空
		Address: address,            // 服务 IP
		Meta:    map[string]string{"Protocol": protocol},
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register server error : ", err)
		return
	}
	return true, nil
}

func setServiceCheckDefaultValue(check *api.AgentServiceCheck) *api.AgentServiceCheck {
	if check.Timeout == "" {
		check.Timeout = "1s"
	}
	if check.Interval == "" {
		// 健康检查间隔
		check.Interval = "1s"
	}
	if check.DeregisterCriticalServiceAfter == "" {
		// check失败后30秒删除本服务，注销时间，相当于过期时间
		check.DeregisterCriticalServiceAfter = "90s"
	}
	return check
}
