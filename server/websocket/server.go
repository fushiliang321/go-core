package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/helper"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type WsServer struct {
	Ctx                   *fasthttp.RequestCtx
	Conn                  *websocket.Conn
	Fd                    uint64
	MessageType           int
	LastResponseTimestamp int64
}

var messageType = 0 //消息类型 0客户端定义 1文本 2二进制

func Start() {
	config := server.Get()

	//消息类型
	switch config.Settings.MessageType {
	case websocket.BinaryMessage:
		messageType = websocket.BinaryMessage
	case websocket.TextMessage:
		messageType = websocket.TextMessage
	}

	if config.Settings.HeartbeatCheckInterval > 0 {
		//心跳检测
		go heartbeatCheck(config.Settings.HeartbeatCheckInterval, config.Settings.HeartbeatIdleTime)
	}
}

func (s *WsServer) Push(data any) {
	bytes, err := helper.AnyToBytes(data)
	if err != nil {
		log.Println("ws push data err:", err)
		return
	}
	err = s.Conn.WriteMessage(s.MessageType, *bytes)
	if err != nil {
		log.Println("ws push write err:", err)
		return
	}
}
func (s *WsServer) Disconnect(data []byte) {
	sender.remove(s.Fd)
	err := s.Conn.WriteControl(websocket.CloseMessage, data, time.Time{})
	if err != nil {
		return
	}
	err = s.Conn.Close()
	if err != nil {
		return
	}
}
func (s *WsServer) Check() bool {
	return Check(s.Fd)
}
