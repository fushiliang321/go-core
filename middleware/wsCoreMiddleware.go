package middleware

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper"
	types2 "github.com/fushiliang321/go-core/router/types"
	"github.com/fushiliang321/go-core/server/types"
	websocket2 "github.com/fushiliang321/go-core/server/websocket"
	"github.com/fushiliang321/go-core/server/websocket/event"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type WebsocketCoreMiddleware struct {
}

var upgraderDefault = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"BinaryMessage", "TextMessage"},
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

var config *server.Servers

func init() {
	config = server.Get()
}

func (m *WebsocketCoreMiddleware) Process(ctx *fasthttp.RequestCtx, handler types2.RequestHandler) (_ any) {
	defer func() {
		exception.Listener("ws process", recover())
	}()
	ws, wsOk := ctx.UserValue(types.SERVER_WEBSOCKET_KEY).(*server.Server)
	if !wsOk {
		ctx.Response.SetStatusCode(500)
		return helper.Error(500, fmt.Sprintln("upgrader exception:"), nil)
	}
	ser := websocket2.SetServer(ctx)
	upgrader := upgraderDefault
	onHandshake, ok := ws.Callbacks[event.ON_HAND_SHAKE].(event.OnHandshake)
	if ok {
		log.Println(ctx.ID(), "onHandshake")
		onHandshake(ser, &upgrader)
	}
	if ctx.Response.StatusCode() != 200 {
		//非200的状态需直接返回，不能升级到websocket
		return
	}
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		ser.Conn = conn
		if ser.MessageType == 0 {
			if conn.Subprotocol() == "BinaryMessage" {
				ser.MessageType = websocket.BinaryMessage
			} else {
				ser.MessageType = websocket.TextMessage
			}
		}
		onOpen, ok := ws.Callbacks[event.ON_OPEN].(event.OnOpen)
		if ok {
			log.Println(ctx.ID(), "onOpen")
			onOpen(ser)
		}

		conn.SetCloseHandler(func(code int, text string) error {
			onClose, ok := ws.Callbacks[event.ON_CLOSE].(event.OnClose)
			if ok {
				log.Println(ctx.ID(), "onClose")
				onClose(ser, code, text)
			}
			return nil
		})

		conn.SetPingHandler(func(appData string) error {
			ser.LastResponseTimestamp = time.Now().Unix()
			conn.WriteControl(websocket.PongMessage, []byte{1}, time.Time{})
			return nil
		})

		if config.Settings.HeartbeatCheckInterval > 0 {
			//心跳检测
			conn.SetPongHandler(func(appData string) error {
				ser.LastResponseTimestamp = time.Now().Unix()
				return nil
			})
		}

		onMessage, ok := ws.Callbacks[event.ON_MESSAGE].(event.OnMessage)
		if ok {
			for {
				_, p, err := conn.ReadMessage()
				if err != nil {
					break
				}
				ser.LastResponseTimestamp = time.Now().Unix()
				func(on event.OnMessage, ser *websocket2.WsServer, p []byte) {
					defer func() {
						if err := recover(); err != nil {
							helper.ErrorResponse(ctx, 500, fmt.Sprintln("ws onMessage exception:", err), nil)
							exception.Listener("ws onMessage exception", err)
						}
					}()
					on(ser, p)
				}(onMessage, ser, p)
			}
		} else {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
				ser.LastResponseTimestamp = time.Now().Unix()
			}
		}
		websocket2.RemoveServer(ser)
	})
	if err != nil {
		websocket2.RemoveServer(ser)
		log.Println("upgrader exception:", err)
		ctx.Response.SetStatusCode(500)
		return helper.Error(500, fmt.Sprintln("upgrader exception:", err), nil)
	}
	return
}
