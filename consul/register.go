package consul

import (
	"fmt"
	"github.com/fushiliang321/go-core/exception"
	"github.com/hashicorp/consul/api"
	"log"
	"strconv"
	"strings"
)

var (
	client               *api.Client
	serviceRegistrations = map[string]*api.AgentServiceRegistration{} //全局的服务注册信息
)

func ServiceRegistrations() *map[string]*api.AgentServiceRegistration {
	return &serviceRegistrations
}

func newConsulClient() (*api.Client, error) {
	var err error
	if client != nil {
		return client, err
	}
	client, err = api.NewClient(GetConfig())
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		client = nil
	}
	return client, err
}

func IsRegister(name string, protocol string, address string, port int) bool {
	_client, _ := newConsulClient()
	_services, err := _client.Agent().Services()
	if err != nil {
		return false
	}
	var service *api.AgentService
	for _, service = range _services {
		if service == nil || service.Service != name || service.Address != address || service.Port != port {
			continue
		}
		MetaProtocol, ok := service.Meta["Protocol"]
		if !ok || MetaProtocol != protocol {
			continue
		}
		return true
	}
	return false
}

func RegisterServer(name string, protocol string, address string, port int, check *api.AgentServiceCheck) (b bool, err error) {
	defer func() {
		exception.Listener("RegisterServer error:", recover())
	}()
	if IsRegister(name, protocol, address, port) {
		return true, nil
	}
	_client, err := newConsulClient()
	if err != nil {
		log.Println("consul client error : ", err)
		return
	}
	if check != nil {
		check = setServiceCheckDefaultValue(check)
	}
	registration := &api.AgentServiceRegistration{
		Name:    name,                               // 服务名称
		Port:    port,                               // 服务端口
		ID:      generateId(getLastServiceId(name)), //服务id
		Tags:    []string{protocol},                 // tag，可以为空
		Address: address,                            // 服务 IP
		Meta:    map[string]string{"Protocol": protocol},
		Check:   check,
	}
	err = _client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register server error : ", err)
		return
	}
	serviceRegistrations[name] = registration
	return true, nil
}

// 设置健康检测默认值
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

// 获取最大的服务id
func getLastServiceId(name string) (maxServiceId string) {
	var (
		err error
		id  int
	)
	maxId := -1
	maxServiceId = name
	_client, _ := newConsulClient()
	_services, err := _client.Agent().Services()
	if err != nil {
		return
	}
	for _, v := range _services {
		if v == nil || v.Service != name {
			continue
		}
		i := strings.LastIndex(v.ID, "-")
		if i == -1 {
			continue
		}
		id, err = strconv.Atoi(v.ID[i+1:])
		if err == nil && id > maxId {
			maxId = id
			maxServiceId = v.ID
		}
	}
	return
}

// 生成id
func generateId(name string) string {
	i := strings.LastIndex(name, "-")
	if i == -1 {
		return name + "-0"
	}
	id, err := strconv.Atoi(name[i+1:])
	name = name[:i]
	if err != nil {
		return name + "-0"
	}
	idStr := strconv.Itoa(id + 1)
	return name + "-" + idStr
}
