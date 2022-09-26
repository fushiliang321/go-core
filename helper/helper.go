package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/fushiliang321/go-core/helper/types"
	"github.com/valyala/fasthttp"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

var appName string

func init() {
	rand.Seed(time.Now().UnixNano())
	appName = GetEnvDefault("APP_NAME", "")
}

func AppName() string {
	return appName
}

// 响应成功数据
func Success(msg string, data any) (res *types.Result) {
	return &types.Result{
		Code:    1,
		Msg:     msg,
		Data:    data,
		ErrCode: 0,
		Service: appName,
	}
}

// 响应错误数据
func Error(errCode int, msg string, data any) (res *types.Result) {
	return &types.Result{
		Code:    0,
		Msg:     msg,
		Data:    data,
		ErrCode: errCode,
		Service: appName,
	}
}

// 响应错误数据
func ErrorResponse(ctx *fasthttp.RequestCtx, errCode int, msg string, data any) {
	ctx.Response.SetStatusCode(errCode)
	res := types.Result{}
	res.Error(errCode, msg, data)
	_, _ = ctx.Write(res.JsonMarshal())
}

// map转struc
func MapToStruc[_type interface {
	int | string
}](m map[_type]any, s any) (err error) {
	arr, err := json.Marshal(m)
	if err != nil {
		return
	}
	return json.Unmarshal(arr, s)
}

// 字符串md5
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

// 取环境变量
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

// 取随机数
func RangeRand(min, max int) int {
	if min > max {
		panic("the min is greater than max!")
	}
	dif := max - min
	return min + rand.Intn(dif)
}

// 转为蛇形字符串
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 当前运行时方法名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	fn := runtime.FuncForPC(pc[0]).Name()
	countSplit := strings.Split(fn, ".")
	if len(countSplit) > 1 {
		return countSplit[1]
	}
	return ""
}

// 当前运行时包名
func RunPackageName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	fn := runtime.FuncForPC(pc[0]).Name()
	countSplit := strings.Split(fn, "/")
	splitLen := len(countSplit)
	if splitLen > 0 {
		countSplit := strings.Split(countSplit[splitLen-1], ".")
		if len(countSplit) > 0 {
			return countSplit[0]
		}
	}
	return ""
}

// json字符串解码
func JsonDecode(str string, v any) error {
	d := json.NewDecoder(bytes.NewReader([]byte(str)))
	d.UseNumber()
	return d.Decode(&v)
}

// 获取当前日期时间
func Time() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

// 获取本机mac地址
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddrs
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

// 获取本机所有ip地址
func GetLocalIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

// 获取本机ip地址
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 获取客户端ip
func ClientIP(ctx *fasthttp.RequestCtx) (string, int) {
	cip := ctx.Request.Header.Peek("client-ip")
	var ip net.IP
	if cip != nil {
		ip = net.ParseIP(string(cip))
	} else {
		ip = ctx.RemoteIP()
	}
	if ip == nil {
		return "", 0
	}
	s := ip.String()
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return s, 4
		case ':':
			return s, 6
		}
	}
	return "", 0
}

// 获取客户端地址（ip+port）
func ClientAddr(ctx *fasthttp.RequestCtx) string {
	ip, v := ClientIP(ctx)
	if ip == "" {
		return ""
	}
	port := ctx.Request.Header.Peek("client-port")
	if port == nil {
		if ip == ctx.RemoteIP().String() {
			return ctx.RemoteAddr().String()
		}
		return ip
	}
	if v == 4 {
		return ip + ":" + string(port)
	}
	return "[" + ip + "]:" + string(port)
}
