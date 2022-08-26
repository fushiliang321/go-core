package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"time"
)

func SetServer(ctx *fasthttp.RequestCtx) (ser *WsServer) {
	ser = &WsServer{}
	ser.Ctx = ctx
	ser.Fd = ctx.ID()
	if messageType != 0 {
		ser.MessageType = messageType
	}
	ser.LastResponseTimestamp = time.Now().Unix()
	sender.add(ser)
	return
}

func RemoveServer(ser *WsServer) {
	sender.remove(ser.Ctx.ID())
	return
}

func Check(fd uint64) (ok bool) {
	_, ok = sender.servers.Load(fd)
	return
}

// 推送消息
func Push(fd uint64, data any) {
	if ser, ok := sender.get(fd); ok {
		ser.Push(data)
	}
}

// 断开连接
func Disconnect(fd uint64, data []byte) {
	if ser, ok := sender.get(fd); ok {
		ser.Disconnect(data)
	}
}

// 心跳检测
func heartbeatCheck(interval int64, idleTime int64) {
	sleep := time.Second * time.Duration(interval)
	if idleTime <= 0 {
		idleTime = interval * 2
	}
	for {
		time.Sleep(sleep)
		nowTime := time.Now().Unix()
		sender.servers.Range(func(key, value any) bool {
			defer func() {
				recover()
			}()
			if ser, ok := value.(*WsServer); ok {
				if nowTime-ser.LastResponseTimestamp > idleTime {
					//超时断开连接
					ser.Disconnect([]byte("timeout"))
				} else {
					ser.Conn.WriteControl(websocket.PingMessage, []byte{1}, time.Time{})
				}
			} else {
				//类型有问题的就删掉
				sender.remove(key.(uint64))
			}
			return true
		})
	}
}
