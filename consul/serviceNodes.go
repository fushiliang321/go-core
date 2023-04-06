package consul

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type (
	ServiceNode struct {
		ServiceName  string
		Address      string
		Port         string
		CheckStatus  string
		Protocol     string
		IsRemove     bool //节点是否被移除
		onRemoveFuns []func()
	}

	serviceNodeLists = []*ServiceNode

	services struct {
		maps sync.Map
	}
)

const (
	HttpProtocol = "http"
	TcpProtocol  = "tcp"
	GrpcProtocol = "grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 判断服务信息是否存在
func (sers *services) hasService(serviceName string) (ok bool) {
	_, ok = sers.maps.Load(serviceName)
	return
}

// 设置服务节点信息
func (sers *services) setServiceNodes(serviceName string, serviceNodes []*ServiceNode) {
	nodeMap := map[string]serviceNodeLists{}
	for i := range serviceNodes {
		node := serviceNodes[i]
		if node.ServiceName != serviceName {
			continue
		}
		switch node.Protocol {
		case HttpProtocol, GrpcProtocol, TcpProtocol:
			if _, ok := nodeMap[node.Protocol]; !ok {
				nodeMap[node.Protocol] = serviceNodeLists{}
			}
			nodeMap[node.Protocol] = append(nodeMap[node.Protocol], node)
		}
	}
	sers.triggerRemoveServiceNode(serviceName)
	sers.maps.Store(serviceName, nodeMap)
}

// 触发服务节点被移除事件
func (sers *services) triggerRemoveServiceNode(serviceName string) {
	v, ok := sers.maps.Load(serviceName)
	if !ok {
		return
	}
	nodeMap := v.(map[string]serviceNodeLists)
	go func() {
		for _, lists := range nodeMap {
			for _, node := range lists {
				node.triggerRemove()
			}
		}
	}()
}

// 随机取一个节点
func (sers *services) getRandomNode(serviceName string, protocol string) (node *ServiceNode, err error) {
	res, ok := sers.maps.Load(serviceName)
	if !ok {
		err = errors.New("没有匹配到服务数据")
		return
	}
	nodeMap := res.(map[string]serviceNodeLists)
	nodeList, ok := nodeMap[protocol]
	if !ok {
		err = errors.New("没有可用协议")
		return
	}
	nodeLen := len(nodeList)
	switch nodeLen {
	case 0:
		err = errors.New("没有可用节点")
		return
	case 1:
		node = nodeList[0]
	default:
		node = nodeList[rand.Intn(nodeLen)]
	}
	return
}

// 触发节点被移除事件
func (node *ServiceNode) triggerRemove() {
	node.IsRemove = true
	if node.onRemoveFuns == nil {
		return
	}
	for _, fun := range node.onRemoveFuns {
		fun()
	}
}

// 监听节点被移除事件
func (node *ServiceNode) OnRemove(_func func()) {
	if node.onRemoveFuns == nil {
		node.onRemoveFuns = []func(){}
	}
	node.onRemoveFuns = append(node.onRemoveFuns, _func)
}
