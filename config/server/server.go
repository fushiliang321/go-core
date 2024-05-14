package server

import "github.com/valyala/fasthttp"

type (
	Callbacks = map[string]any
	Settings  struct {
		HeartbeatCheckInterval int64 //心跳检测间隔
		HeartbeatIdleTime      int64 //心跳超时时间
		MessageType            int   //websocket消息类型
		AutoResponseGzipSize   int   //响应数据达到指定大小自动触发gzip压缩
		TLS                    *TLS  //tls配置
	}
	TLS struct {
		CertFile string
		KeyFile  string
	}
	Server struct {
		Name      string
		Type      byte
		Host      string
		Port      string
		Callbacks Callbacks
		TLS       *TLS //每个服务可以单独配置tls
		Server    *fasthttp.Server
	}
	Servers struct {
		Servers  []Server
		Settings *Settings
	}
)

var servers = &Servers{
	Servers:  []Server{},
	Settings: &Settings{},
}

func Set(s *Servers) {
	servers = s
}

func Get() *Servers {
	return servers
}
