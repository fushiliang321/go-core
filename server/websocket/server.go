package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type ConnWriteChanParams struct {
	messageType int
	data        *[]byte
	deadline    *time.Time
}

type WsServer struct {
	Ctx                   *fasthttp.RequestCtx
	Conn                  *websocket.Conn
	ConnWriteChan         chan *ConnWriteChanParams
	Fd                    uint64
	MessageType           int
	LastResponseTimestamp int64
	Status                byte
}

const WsServerStatusClose = 0 //连接关闭
const WsServerStatusOpen = 1  //连接开启

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

func (s *WsServer) init() {
	s.ConnWriteChan = make(chan *ConnWriteChanParams, 1)
	go func() {
		defer func() {
			exception.Listener("ws dispose exception", recover())
		}()
		var (
			err       error
			writeData *ConnWriteChanParams
		)
		for writeData = range s.ConnWriteChan {
			if writeData.messageType == s.MessageType {
				if err = s.Conn.WriteMessage(s.MessageType, *writeData.data); err != nil {
					log.Println("ws write message err:", err)
				}
			} else {
				if err = s.Conn.WriteControl(writeData.messageType, *writeData.data, *writeData.deadline); err != nil {
					log.Println("ws write control err:", err)
				}
				if writeData.messageType == websocket.CloseMessage {
					s.Conn.Close()
				}
			}
		}
		close(s.ConnWriteChan)
	}()
}

func (s *WsServer) Push(data any) {
	if s.Status != WsServerStatusOpen {
		return
	}
	bytes, err := helper.AnyToBytes(data)
	if err != nil {
		log.Println("ws push data err:", err)
		return
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: s.MessageType,
		data:        bytes,
	}
}

func (s *WsServer) Ping(data []byte, deadline time.Time) {
	if s.Status != WsServerStatusOpen {
		return
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.PingMessage,
		data:        &data,
		deadline:    &deadline,
	}
}

func (s *WsServer) Pong(data []byte, deadline time.Time) {
	if s.Status != WsServerStatusOpen {
		return
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.PongMessage,
		data:        &data,
		deadline:    &deadline,
	}
}

func (s *WsServer) Disconnect(data []byte) {
	if s.Status != WsServerStatusOpen {
		return
	}
	sender.remove(s.Fd)
	s.Status = WsServerStatusClose
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.CloseMessage,
		data:        &data,
		deadline:    &time.Time{},
	}
}

func (s *WsServer) Check() bool {
	return Check(s.Fd)
}
