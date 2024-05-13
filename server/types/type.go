package types

const (
	SERVER_WEBSOCKET = byte(0) //websocket
	SERVER_HTTP      = byte(1) //http1
	SERVER_HTTP2     = byte(2) //http2
	SERVER_HTTP3     = byte(3) //http3

	SERVER_HTTP_KEY      = "SERVER_KEY/1"
	SERVER_WEBSOCKET_KEY = "SERVER_KEY/2"
	SERVER_NAME_KEY      = "SERVER_NAME_KEY/0"
)
