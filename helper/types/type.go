package types

import (
	"bytes"
	"encoding/json"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/fushiliang321/go-core/helper/system"
)

type Result struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	ErrCode int    `json:"errCode"`
	Service string `json:"service"`
}
type WsResult struct {
	Result
	Path string `json:"path"`
	Mark string `json:"mark,omitempty"`
}

var appName string

func init() {
	appName = system.AppName()
}

func (res *Result) Error(errCode int, msg string, data any) {
	res.Code = 0
	res.ErrCode = errCode
	res.Msg = msg
	if data == nil {
		res.Data = map[string]string{}
	} else {
		res.Data = data
	}
	res.Service = appName
}

// 转成json
func (res *Result) JsonMarshal() (marshal []byte) {
	marshal, err := json.Marshal(res)
	if err != nil {
		logger.Warn("server result err:", err)
		return
	}
	return
}

// 转为指定类型数据
func (res *Result) To(_type any) {
	marshal, err := json.Marshal(res)
	if err != nil {
		return
	}
	d := json.NewDecoder(bytes.NewReader(marshal))
	d.UseNumber()
	err = d.Decode(_type)
	if err != nil {
		return
	}
	return
}
