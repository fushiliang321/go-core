package consul

import (
	"fmt"
	"github.com/fushiliang321/go-core/event"
	"github.com/hashicorp/consul/api"
	"sync"
)

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
	event.Dispatch(event.NewRegistered(event.ConsulConsumerServerStart, serviceName))
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
