package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fushiliang321/go-core/helper/types"
	"github.com/fushiliang321/go-core/jsonRpcHttp/context"
	"github.com/fushiliang321/go-core/rpc"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var appName string
var cacheData = map[string]any{}

func init() {
	rand.Seed(time.Now().UnixNano())
	appName = GetEnvDefault("APP_NAME", "")
}

func AppName() string {
	return appName
}

// 响应成功数据
func Success(msg string, datas ...any) (res *types.Result) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	return &types.Result{
		Code:    1,
		Msg:     msg,
		Data:    data,
		ErrCode: 0,
		Service: appName,
	}
}

// 响应错误数据
func Error(errCode int, msg string, datas ...any) (res *types.Result) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	return &types.Result{
		Code:    0,
		Msg:     msg,
		Data:    data,
		ErrCode: errCode,
		Service: appName,
	}
}

// 响应错误数据
func ErrorResponse(ctx *fasthttp.RequestCtx, errCode int, msg string, datas ...any) {
	var data any
	if len(datas) == 0 || datas[0] == nil {
		data = map[string]string{}
	} else {
		data = datas[0]
	}
	ctx.Response.SetStatusCode(errCode)
	res := types.Result{}
	res.Error(errCode, msg, data)
	_, _ = ctx.Write(res.JsonMarshal())
}

// map转struc
func MapToStruc[_type interface {
	int | string
}](m map[_type]any, s any) (err error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		return
	}
	d := json.NewDecoder(bytes.NewReader(marshal))
	d.UseNumber()
	return d.Decode(s)
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

// json编码
func JsonEncode(v any) (string, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
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

// 获取本机ip地址，默认获取对外的ip地址
func GetLocalIP(address ...string) string {
	var err error
	address = append(address, []string{
		"223.5.5.5:53", "8.8.8.8:53",
	}...)
	for _, addr := range address {
		ip := net.ParseIP(addr)
		if ip != nil {
			switch IpType(ip.String()) {
			case 4:
				addr = (addr + ":80")
			case 6:
				addr = ("[" + addr + "]:80")
			}
		}
		conn, err := net.Dial("udp", addr)
		if err != nil {
			continue
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}
	log.Println("GetLocalIP error：", err)
	return ""
}

// 判断ip类型
func IpType(ip string) uint8 {
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	return 0
}

// 获取客户端ip
func ClientIP(ctx *fasthttp.RequestCtx) (string, uint8) {
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
	return s, IpType(s)
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

// 获取结构体字段
func GetStructFields(v any) (fields []string) {
	if reflect.ValueOf(v).Kind() == reflect.Struct {
		reflectType := reflect.TypeOf(v)
		numField := reflectType.NumField()
		for i := 0; i < numField; i++ {
			fields = append(fields, reflectType.Field(i).Name)
		}
	}
	return
}

// 获取文件内的所有变量名
func GetFileVariateNameAll(FilePath string) (names []string) {
	cache, ok := cacheData["GetFileVariateNameAll"]
	mapData := map[string][]string{}
	if ok {
		if mapData1, ok := cache.(map[string][]string); ok {
			mapData = mapData1
			if names, ok := mapData[FilePath]; ok {
				return names
			}
		}
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, FilePath, nil, 0)
	if err != nil {
		fmt.Printf("err = %s", err)
	}
	for s, object := range f.Scope.Objects {
		if object.Kind == ast.Var {
			names = append(names, s)
		}
	}
	mapData[FilePath] = names
	cacheData["GetFileVariateNameAll"] = mapData
	return
}

// 获取当前文件名
//
//go:noinline
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Can not get current file info")
	}
	return file
}

// 获取rpc上下文请求数据
func RpcRequestData() (rpcRequestData rpc.RpcRequestData) {
	ctxData := context.GetAll()
	if ctxData == nil {
		return
	}
	data := ctxData["internalRequest"]
	if data == nil {
		return
	}
	mapData, ok := data.(map[string]any)
	if !ok {
		return
	}
	err := MapToStruc[string](mapData, &rpcRequestData)
	if err != nil {
		return
	}
	return rpcRequestData
}

func AnyToBytes(data any) (*[]byte, error) {
	var bts []byte
	var err error
	switch data.(type) {
	case string:
		bts = strconv.S2B(data.(string))
	case *string:
		bts = strconv.S2B(*(data.(*string)))
	case []byte:
		bts = data.([]byte)
	case *[]byte:
		return data.(*[]byte), nil
	case byte:
		bts = []byte{data.(byte)}
	case *byte:
		bts = []byte{*(data.(*byte))}
	default:
		bts, err = json.Marshal(data)
	}
	return &bts, err
}
