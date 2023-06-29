package server

import "github.com/valyala/fasthttp"

type (
	Callbacks = map[string]any
	Settings  struct {
		HeartbeatCheckInterval int64 //心跳检测间隔
		HeartbeatIdleTime      int64 //心跳超时时间
		MessageType            int   //websocket消息类型
	}
	Server struct {
		Name      string
		Type      byte
		Host      string
		Port      string
		Callbacks Callbacks
		Server    *fasthttp.Server
	}
	Servers struct {
		Servers  []Server
		Settings *Settings
	}
)

var servers = &Servers{
	Settings: &Settings{},
	Servers:  []Server{},
}

func Set(s *Servers) {
	servers = s
}

func Get() *Servers {
	return servers
}
