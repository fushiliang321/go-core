package consul

import (
	"fmt"
	"github.com/fushiliang321/go-core/event"
	"github.com/fushiliang321/go-core/helper/logger"
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
			logger.Error("error getService", fmt.Sprint(err))
		}
		if s.status == monitorOff {
			return
		}
		go func() {
			time.Sleep(time.Second * 1)
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
		logger.Warn("error retrieving instances from Consul: ", err)
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
		port         string
	)
	for _, ser = range sers {
		port = strconv.Itoa(ser.Service.Port)
		for _, check = range ser.Checks {
			if check.Status == api.HealthPassing && check.Type != "" {
				serviceNodes = append(serviceNodes, &ServiceNode{
					CheckStatus: api.HealthPassing,
					ServiceName: check.ServiceName,
					Protocol:    check.Type,
					Address:     ser.Service.Address,
					Port:        port,
					IsRemove:    false,
				})
			}
		}
	}
	if serviceNodes == nil {
		serviceNodes = []*ServiceNode{}
	}
	globalServices.setServiceNodes(s.name, serviceNodes)
	event.Dispatch(event.NewRegistered(event.ConsulConsumerServiceInfoChange, s.name))
}

func (s *serviceMonitor) close() {
	s.status = monitorOff
}
