package websocket

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/router/types"
	"time"
)

type (
	ConnWriteChanParams struct {
		messageType int
		data        []byte
		deadline    time.Time
	}
	WsServer struct {
		Ctx                   *types.RequestCtx
		Conn                  *websocket.Conn
		ConnWriteChan         chan *ConnWriteChanParams
		Fd                    uint64
		MessageType           int
		LastResponseTimestamp int64
		Status                byte
	}
)

const (
	WsServerStatusClose     = 0 //连接关闭
	WsServerStatusOpen      = 1 //连接开启
	WsServerStatusBeClosing = 2 //连接正在关闭
)

var (
	DataFramesDefault []byte
	messageType       = 0 //消息类型 0客户端定义 1文本 2二进制
)

func Start() {
	config := server.Get()

	if config.Settings == nil {
		return
	}

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
			if err := recover(); err != nil {
				logger.Error("ws dispose exception:", fmt.Sprint(err))
				exception.Listener("ws dispose exception", err)
			}
		}()
		var (
			err       error
			writeData *ConnWriteChanParams
		)
		for writeData = range s.ConnWriteChan {
			if s.Status == WsServerStatusClose {
				return
			}
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						s.Conn.Close()
						logger.Error("["+fmt.Sprint(s.Fd)+"]ws write message exception", fmt.Sprint(err), writeData.messageType, writeData.data, writeData.deadline, s.Conn.NetConn())
						exception.Listener("["+fmt.Sprint(s.Fd)+"]ws write message exception", rec)
					}
				}()
				switch writeData.messageType {
				case s.MessageType: //发送消息帧
					if err = s.Conn.WriteMessage(s.MessageType, writeData.data); err != nil {
						logger.Warn("ws write message err:", err)
					}
				case websocket.CloseMessage: //发送关闭帧
					s.Conn.WriteControl(writeData.messageType, writeData.data, writeData.deadline)
					s.Conn.Close()
				default: //发送其他控制帧
					if err = s.Conn.WriteControl(writeData.messageType, writeData.data, writeData.deadline); err != nil {
						logger.Warn("ws write control err:", err)
					}
				}
			}()
		}
	}()
}

func (s *WsServer) Push(data any) {
	if s.Status != WsServerStatusOpen {
		return
	}
	bytes, err := helper.AnyToBytes(data)
	if err != nil {
		logger.Warn("ws push data err:", err)
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
	if data == nil {
		data = DataFramesDefault
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.PingMessage,
		data:        data,
		deadline:    deadline,
	}
}

func (s *WsServer) Pong(data []byte, deadline time.Time) {
	if s.Status != WsServerStatusOpen {
		return
	}
	if data == nil {
		data = DataFramesDefault
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.PongMessage,
		data:        data,
		deadline:    deadline,
	}
}

// 连接断开事件处理
func (s *WsServer) OnClose() {
	if s.Status == WsServerStatusClose {
		return
	}
	sender.remove(s.Fd)
	s.Status = WsServerStatusClose
	close(s.ConnWriteChan)
}

// 主动断开连接
func (s *WsServer) Disconnect(data []byte) {
	if s.Status != WsServerStatusOpen {
		return
	}
	sender.remove(s.Fd)
	s.Status = WsServerStatusBeClosing
	if data == nil {
		data = DataFramesDefault
	}
	s.ConnWriteChan <- &ConnWriteChanParams{
		messageType: websocket.CloseMessage,
		data:        data,
		deadline:    time.Time{},
	}
}

func (s *WsServer) Check() bool {
	return Check(s.Fd)
}
