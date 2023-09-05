package middleware

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/response"
	types2 "github.com/fushiliang321/go-core/router/types"
	"github.com/fushiliang321/go-core/server/types"
	websocket2 "github.com/fushiliang321/go-core/server/websocket"
	"github.com/fushiliang321/go-core/server/websocket/event"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
	"time"
)

type (
	WebsocketCoreMiddleware struct{}
	wsError                 struct {
		code int
		text string
	}
)

var (
	upgraderDefault = websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"BinaryMessage", "TextMessage"},
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return true
		},
	}
	config *server.Servers

	pongData = []byte{1}
)

func init() {
	config = server.Get()
}

func (m *WebsocketCoreMiddleware) Process(ctx *types2.RequestCtx, handler types2.RequestHandler) (_ any) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("ws process error:", fmt.Sprint(err))
			exception.Listener("ws process", err)
		}
	}()
	ws, ok := ctx.Raw().UserValue(types.SERVER_WEBSOCKET_KEY).(*server.Server)
	if !ok {
		ctx.Response.SetStatusCode(500)
		return response.Error(500, fmt.Sprintln("upgrader exception:"), nil)
	}
	ser := websocket2.SetServer(ctx)
	upgrader := upgraderDefault
	onHandshake, ok := ws.Callbacks[event.ON_HAND_SHAKE].(event.OnHandshake)
	if ok {
		logger.Debug(ctx.Raw().ID(), "onHandshake")
		onHandshake(ser, &upgrader)
	}
	if ctx.Response.StatusCode() != 200 {
		//非200的状态需直接返回，不能升级到websocket
		return
	}
	err := upgrader.Upgrade(ctx.Raw(), func(conn *websocket.Conn) {
		ser.Conn = conn

		do := sync.Once{}
		closeHandler := func(code int, text string) error {
			defer func() {
				if err := recover(); err != nil {
					logger.Warn("ws on close handler error:", fmt.Sprint(err))
					exception.Listener("ws on close handler", err)
				}
			}()

			do.Do(func() {
				conn.Close()
				ser.OnClose()
				if onClose, ok := ws.Callbacks[event.ON_CLOSE].(event.OnClose); ok {
					logger.Debug(ctx.Raw().ID(), "onClose")
					onClose(ser, code, text)
				}
			})
			return nil
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Error("ws handler error:", fmt.Sprint(err))
				exception.Listener("ws handler", err)
			}
			closeHandler(0, "")
		}()

		if ser.MessageType == 0 {
			if conn.Subprotocol() == "BinaryMessage" {
				ser.MessageType = websocket.BinaryMessage
			} else {
				ser.MessageType = websocket.TextMessage
			}
		}

		if onOpen, ok := ws.Callbacks[event.ON_OPEN].(event.OnOpen); ok {
			logger.Debug(ctx.Raw().ID(), "onOpen")
			onOpen(ser)
		}

		conn.SetCloseHandler(closeHandler)

		conn.SetPingHandler(func(appData string) error {
			//响应ping帧
			ser.LastResponseTimestamp = time.Now().Unix()
			ser.Pong(pongData, time.Time{})
			return nil
		})

		if config.Settings.HeartbeatCheckInterval > 0 {
			//心跳检测
			conn.SetPongHandler(func(appData string) error {
				ser.LastResponseTimestamp = time.Now().Unix()
				return nil
			})
		}
		var err error
		if onMessage, ok := ws.Callbacks[event.ON_MESSAGE].(event.OnMessage); ok {
			var p []byte
			for {
				_, p, err = conn.ReadMessage()
				if err != nil {
					_wsError := errAnalysis(err)
					closeHandler(_wsError.code, _wsError.text)
					break
				}
				ser.LastResponseTimestamp = time.Now().Unix()
				callOnMessage(onMessage, ser, p)
				p = nil
			}
		} else {
			for {
				if _, _, err = conn.ReadMessage(); err != nil {
					_wsError := errAnalysis(err)
					closeHandler(_wsError.code, _wsError.text)
					break
				}
				ser.LastResponseTimestamp = time.Now().Unix()
			}
		}
	})
	if err != nil {
		websocket2.RemoveServer(ser)
		ctx.Response.SetStatusCode(500)
		return response.Error(500, fmt.Sprintln("upgrader exception:", err), nil)
	}
	return
}

// 调用OnMessage方法
func callOnMessage(on event.OnMessage, ser *websocket2.WsServer, p []byte) (err any) {
	defer func() {
		if err = recover(); err != nil {
			ser.Push(response.Error(500, fmt.Sprintln("ws onMessage exception:", err)))
			logger.Error("ws onMessage exception:", fmt.Sprint(err))
			exception.Listener("ws onMessage exception", err)
		}
	}()
	on(ser, p)
	return
}

func errAnalysis(_err error) *wsError {
	code := 0
	str := _err.Error()
	if len(str) > 20 {
		str = str[17:21]
		i, err := strconv.Atoi(str)
		if err == nil {
			if strconv.Itoa(i) == str {
				code = i
			}
		}
	}
	return &wsError{
		code: code,
		text: _err.Error(),
	}
}
