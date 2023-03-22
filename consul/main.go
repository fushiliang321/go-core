package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
)

type Service struct{}

var (
	apiConfig *api.Config
	services  = Services{}

	serviceMap     map[string]*serviceMonitor
	serviceMapLock sync.Locker
)

func newClient() (client *api.Client, err error) {
	client, err = api.NewClient(apiConfig)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
	}
	return
}

func AddServices(consumerServiceNames []string) {
	for _, serviceName := range consumerServiceNames {
		go AddService(serviceName)
	}
}

// 添加服务信息
func AddService(serviceName string) {
	serviceMapLock.Lock()
	defer serviceMapLock.Unlock()
	if _, ok := serviceMap[serviceName]; ok {
		return
	}
	serviceMap[serviceName] = &serviceMonitor{
		name:      serviceName,
		status:    monitorOn,
		lastIndex: uint64(0),
	}
	serviceMap[serviceName].syncService()
}

// 移除服务信息
func RemoveService(serviceName string) {
	serviceMapLock.Lock()
	defer serviceMapLock.Unlock()
	serviceMap[serviceName].close()
	delete(serviceMap, serviceName)
}

func GetNode(serviceName string, protocol string) (node *ServiceNode, err error) {
	return services.getRandomNode(serviceName, protocol)
}
