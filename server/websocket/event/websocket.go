package event

import (
	websocket2 "core/server/websocket"
	"github.com/fasthttp/websocket"
)

type (
	OnHandshake = func(ser *websocket2.WsServer, upgrader *websocket.FastHTTPUpgrader)
	OnOpen      = func(ser *websocket2.WsServer)
	OnMessage   = func(ser *websocket2.WsServer, data []byte)
	OnClose     = func(ser *websocket2.WsServer, code int, text string)
	Websocket   interface {
		OnHandshake(ser *websocket2.WsServer, upgrader *websocket.FastHTTPUpgrader)
		OnOpen(ser *websocket2.WsServer)
		OnMessage(ser *websocket2.WsServer, data []byte)
		OnClose(ser *websocket2.WsServer, code int, text string)
	}
)

const (
	ON_HAND_SHAKE = "handshake"
	ON_OPEN       = "open"
	ON_MESSAGE    = "message"
	ON_CLOSE      = "close"
)
