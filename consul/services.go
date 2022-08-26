package consul

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type ServiceNode struct {
	ServiceName string
	Address     string
	Port        string
	Status      string
	Protocol    string
}
type ServiceNodes struct {
	nodes []ServiceNode
}

type Services struct {
	maps sync.Map
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (sers *Services) isExistService(serviceName string) (ok bool) {
	_, ok = sers.maps.Load(serviceName)
	return
}

func (sers *Services) setServiceNodes(serviceName string, serviceNodes []ServiceNode) {
	nodeMap := map[string]ServiceNodes{}
	tcpNodes := ServiceNodes{}
	httpNodes := ServiceNodes{}
	for _, node := range serviceNodes {
		if node.ServiceName != serviceName {
			continue
		}
		if node.Protocol == "tcp" {
			tcpNodes.nodes = append(tcpNodes.nodes, node)
		} else if node.Protocol == "http" {
			httpNodes.nodes = append(httpNodes.nodes, node)
		}
	}
	nodeMap["tcp"] = tcpNodes
	nodeMap["http"] = httpNodes

	sers.maps.Store(serviceName, nodeMap)
}

func (sers *Services) getRandomNode(serviceName string) (node ServiceNode, err error) {
	res, ok := sers.maps.Load(serviceName)
	if !ok {
		err = errors.New("没有匹配到服务数据")
		return
	}
	nodeMap := res.(map[string]ServiceNodes)
	nodeLen := len(nodeMap["http"].nodes)
	tcpNodeLen := len(nodeMap["tcp"].nodes)
	if tcpNodeLen == 0 && nodeLen == 0 {
		err = errors.New("没有可用节点")
		return
	}
	protocol := "http"
	if tcpNodeLen > 0 && rand.Intn(2) == 0 {
		protocol = "tcp"
		nodeLen = tcpNodeLen
	}
	node = nodeMap[protocol].nodes[rand.Intn(nodeLen)]
	return
}
