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
	nodes []*ServiceNode
}

type Services struct {
	maps sync.Map
}

var (
	HttpProtocol = "http"
	TcpProtocol  = "tcp"
	GrpcProtocol = "grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (sers *Services) isExistService(serviceName string) (ok bool) {
	_, ok = sers.maps.Load(serviceName)
	return
}

func (sers *Services) setServiceNodes(serviceName string, serviceNodes []*ServiceNode) {
	nodeMap := map[string]ServiceNodes{}
	tcpNodes := ServiceNodes{}
	httpNodes := ServiceNodes{}
	grpcNodes := ServiceNodes{}
	for _, node := range serviceNodes {
		if node.ServiceName != serviceName {
			continue
		}
		switch node.Protocol {
		case TcpProtocol:
			tcpNodes.nodes = append(tcpNodes.nodes, node)
		case HttpProtocol:
			httpNodes.nodes = append(httpNodes.nodes, node)
		case GrpcProtocol:
			grpcNodes.nodes = append(grpcNodes.nodes, node)
		}
	}
	nodeMap[TcpProtocol] = tcpNodes
	nodeMap[HttpProtocol] = httpNodes
	nodeMap[GrpcProtocol] = grpcNodes

	sers.maps.Store(serviceName, nodeMap)
}

func (sers *Services) getRandomNode(serviceName string, protocol string) (node *ServiceNode, err error) {
	res, ok := sers.maps.Load(serviceName)
	if !ok {
		err = errors.New("没有匹配到服务数据")
		return
	}
	nodeMap := res.(map[string]ServiceNodes)
	nodeInfo, ok := nodeMap[protocol]
	if !ok {
		err = errors.New("没有可用协议")
		return
	}
	nodeLen := len(nodeInfo.nodes)
	if nodeLen == 0 {
		err = errors.New("没有可用节点")
		return
	}
	node = nodeInfo.nodes[rand.Intn(nodeLen)]
	return
}
