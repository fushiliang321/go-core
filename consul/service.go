package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"time"
)

type serviceMonitor struct {
	name      string //服务名称
	status    byte   //状态0结束监听 1开始监听
	lastIndex uint64
}

const (
	monitorOff = 0
	monitorOn  = 1
)

// 同步服务信息
func (s *serviceMonitor) syncService() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error getService", err)
		}
		if s.status == monitorOff {
			return
		}
		go func() {
			time.Sleep(time.Millisecond * 10)
			s.syncService()
		}()
	}()

	if s.status == monitorOff {
		return
	}
	_client, _ := newClient()
	sers, metaInfo, err := _client.Health().Service(s.name, "", true, &api.QueryOptions{
		WaitIndex: s.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		fmt.Println("error retrieving instances from Consul: ", err)
		return
	}
	if s.status == monitorOff {
		return
	}
	s.lastIndex = metaInfo.LastIndex
	var (
		serviceNodes []*ServiceNode
		ser          *api.ServiceEntry
		check        *api.HealthCheck
	)
	for _, ser = range sers {
		for _, check = range ser.Checks {
			if check.Status == api.HealthPassing && check.Type != "" {
				node := &ServiceNode{}
				node.Status = check.Status
				node.ServiceName = check.ServiceName
				node.Protocol = check.Type
				node.Address = ser.Service.Address
				node.Port = strconv.Itoa(ser.Service.Port)
				serviceNodes = append(serviceNodes, node)
			}
		}
	}
	if serviceNodes == nil {
		serviceNodes = []*ServiceNode{}
	}
	services.setServiceNodes(s.name, serviceNodes)
}

func (s *serviceMonitor) close() {
	s.status = monitorOff
}
