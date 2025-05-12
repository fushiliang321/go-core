package server

import "github.com/valyala/fasthttp"

const (
	//压缩类型
	CompressTypeAuto    CompressType = iota //自动
	CompressTypeGzip                        //gzip
	CompressTypeDeflate                     //deflate
	CompressTypeBrotli                      //brotli
	CompressTypeZstd                        //zstd
)

type (
	CompressType = uint8
	Callbacks    = map[string]any
	Settings     struct {
		HeartbeatCheckInterval int64     //心跳检测间隔
		HeartbeatIdleTime      int64     //心跳超时时间
		MessageType            int       //websocket消息类型
		Compress               *Compress //数据压缩配置
		TLS                    *TLS      //tls配置
	}
	Compress struct {
		Type    CompressType //压缩类型
		MinSize int          //响应数据超过指定大小时触发压缩，0为全部压缩
		Level   int          //压缩等级，0为默认
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
