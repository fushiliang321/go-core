package websocket

import (
	"sync"
)

type Sender struct {
	servers sync.Map
}

var sender = Sender{}

func (sender *Sender) add(ser *WsServer) {
	sender.servers.Store(ser.Fd, ser)
}

func (sender *Sender) get(fd uint64) (s *WsServer, o bool) {
	if ser, ok := sender.servers.Load(fd); ok {
		if s, o = ser.(*WsServer); !o {
			//类型有问题的就删掉
			sender.remove(fd)
		}
	}
	return
}

func (sender *Sender) remove(fd uint64) {
	sender.servers.Delete(fd)
}
