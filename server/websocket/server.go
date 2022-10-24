package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
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
	var msgData []byte
	switch data.(type) {
	case byte:
		msgData = []byte{data.(byte)}
	case *byte:
		msgData = []byte{*(data.(*byte))}
	case []byte:
		msgData = data.([]byte)
	case *[]byte:
		msgData = *(data.(*[]byte))
	case string:
		msgData = []byte(data.(string))
	case *string:
		msgData = []byte(*(data.(*string)))
	default:
		err := s.Conn.WriteJSON(data)
		if err != nil {
			log.Println("ws push err:", err)
		}
		return
	}
	err := s.Conn.WriteMessage(s.MessageType, msgData)
	if err != nil {
		log.Println("ws push err:", err)
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
