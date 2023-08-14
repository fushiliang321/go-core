package websocket

import (
	"github.com/fushiliang321/go-core/router/types"
	"github.com/savsgio/gotils/strconv"
	"time"
)

func SetServer(ctx *types.RequestCtx) (ser *WsServer) {
	ser = &WsServer{
		Ctx:                   ctx,
		Fd:                    ctx.Raw().ID(),
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
	sender.remove(ser.Ctx.Raw().ID())
	return
}

func Check(fd uint64) (ok bool) {
	_, ok = sender.Load(fd)
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
		nowTime        int64
		ok             bool
		ser            *WsServer
		pingData       = []byte{1}
		disconnectData = strconv.S2B("timeout")
	)
	for {
		time.Sleep(sleep)
		nowTime = time.Now().Unix()
		sender.Range(func(fd, value any) bool {
			defer func() {
				recover()
			}()
			if ser, ok = value.(*WsServer); ok && ser.Status == WsServerStatusOpen {
				if nowTime-ser.LastResponseTimestamp > idleTime {
					//超时断开连接
					ser.Disconnect(disconnectData)
				} else {
					ser.Ping(pingData, time.Time{})
				}
			} else {
				//类型或者连接状态有问题的就删掉
				sender.remove(fd.(uint64))
			}
			return true
		})
		ser = nil
	}
}
