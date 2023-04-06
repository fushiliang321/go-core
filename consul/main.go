package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
)

type Service struct{}

var (
	apiConfig      *api.Config
	globalServices = services{}

	serviceMonitorMap     = map[string]*serviceMonitor{}
	serviceMonitorMapLock sync.Mutex
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
	serviceMonitorMapLock.Lock()
	defer serviceMonitorMapLock.Unlock()
	if _, ok := serviceMonitorMap[serviceName]; ok {
		return
	}
	serviceMonitorMap[serviceName] = &serviceMonitor{
		name:      serviceName,
		status:    monitorOn,
		lastIndex: uint64(0),
	}
	serviceMonitorMap[serviceName].syncService()
}

// 移除服务信息
func RemoveService(serviceName string) {
	serviceMonitorMapLock.Lock()
	defer serviceMonitorMapLock.Unlock()
	serviceMonitorMap[serviceName].close()
	delete(serviceMonitorMap, serviceName)
}

func GetNode(serviceName string, protocol string) (node *ServiceNode, err error) {
	return globalServices.getRandomNode(serviceName, protocol)
}
