package websocket

import (
	"github.com/valyala/fasthttp"
	"time"
)

func SetServer(ctx *fasthttp.RequestCtx) (ser *WsServer) {
	ser = &WsServer{
		Ctx:                   ctx,
		Fd:                    ctx.ID(),
		Status:                WsServerStatusOpen,
		LastResponseTimestamp: time.Now().Unix(),
	}
	if messageType != 0 {
		ser.MessageType = messageType
	}
	ser.init()
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
	var (
		nowTime int64
		ok      bool
		ser     *WsServer
	)
	for {
		time.Sleep(sleep)
		nowTime = time.Now().Unix()
		sender.servers.Range(func(fd, value any) bool {
			defer func() {
				recover()
			}()
			if ser, ok = value.(*WsServer); ok {
				if nowTime-ser.LastResponseTimestamp > idleTime {
					//超时断开连接
					ser.Disconnect([]byte("timeout"))
				} else {
					ser.Ping([]byte{1}, DeadlineDefault)
				}
			} else {
				//类型有问题的就删掉
				sender.remove(fd.(uint64))
			}
			return true
		})
	}
}
