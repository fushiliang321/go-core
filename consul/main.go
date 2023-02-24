package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"sync"
	"time"
)

type Service struct{}

var apiConfig *api.Config
var services = Services{}
var lastIndexMap sync.Map
var lastIndexDefault uint64

func newClient() (client *api.Client, err error) {
	client, err = api.NewClient(apiConfig)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
	}
	return
}

func getServiceData(consumerServiceNames []string) {
	for _, serviceName := range consumerServiceNames {
		go SyncService(serviceName)
	}
}

// 同步服务信息
func SyncService(serviceName string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error getService", err)
		}
		go func(serviceName string) {
			time.Sleep(time.Millisecond * 10)
			SyncService(serviceName)
		}(serviceName)
	}()
	lastIndex, _ := lastIndexMap.LoadOrStore(serviceName, lastIndexDefault)
	_client, _ := newClient()
	sers, metaInfo, err := _client.Health().Service(serviceName, "", true, &api.QueryOptions{
		WaitIndex: lastIndex.(uint64), // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		fmt.Println("error retrieving instances from Consul: ", err)
		return
	}
	lastIndexMap.Store(serviceName, metaInfo.LastIndex)
	serviceNodes := []ServiceNode{}
	for _, ser := range sers {
		for _, check := range ser.Checks {
			if check.Status == "passing" && check.Type != "" {
				node := new(ServiceNode)
				node.Status = check.Status
				node.ServiceName = check.ServiceName
				node.Protocol = check.Type
				node.Address = ser.Service.Address
				node.Port = strconv.Itoa(ser.Service.Port)
				serviceNodes = append(serviceNodes, *node)
			}
		}
	}
	services.setServiceNodes(serviceName, serviceNodes)
}

func GetNode(serviceName string, protocol string) (node ServiceNode, err error) {
	return services.getRandomNode(serviceName, protocol)
}
