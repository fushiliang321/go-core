package types

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fushiliang321/go-core/config/server"
	"github.com/fushiliang321/go-core/event/handles/core"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/helper/serialize"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/url"
	"reflect"
	strconv2 "strconv"
	"strings"
	"time"
)

type RequestCtx fasthttp.RequestCtx

var AutoResponseGzipSize int //响应数据达到指定大小自动触发gzip压缩

func init() {
	go core.AwaitStartFinish(func() {
		serversConfig := server.Get()
		if serversConfig.Settings != nil {
			AutoResponseGzipSize = serversConfig.Settings.AutoResponseGzipSize
		}
	})
}

func (ctx *RequestCtx) Raw() *fasthttp.RequestCtx {
	return (*fasthttp.RequestCtx)(ctx)
}

func (ctx *RequestCtx) WriteAny(data any) (int, error) {
	bytes, err := helper.AnyToBytes(data)
	if err != nil {
		log.Printf("write data err:%s\n", err)
		return 0, err
	}
	if AutoResponseGzipSize > 0 && len(bytes) >= AutoResponseGzipSize {
		ctx.Response.Header.Add("Content-Encoding", "gzip")
		return (*fasthttp.RequestCtx)(ctx).Write(fasthttp.AppendGzipBytes([]byte{}, bytes))
	}
	return (*fasthttp.RequestCtx)(ctx).Write(bytes)
}

func (ctx *RequestCtx) initParams() {
	var (
		body        = (*fasthttp.RequestCtx)(ctx).PostBody()
		paramsStr   = (*fasthttp.RequestCtx)(ctx).QueryArgs().QueryString()
		formDataStr = (*fasthttp.RequestCtx)(ctx).PostArgs().QueryString()
	)
	ctx.SetInputs(splitParams(paramsStr, formDataStr, body))
}

// 拆分字符串参数
func splitParams(queryArgs []byte, formDataStr []byte, body []byte) (params map[string]any) {
	var (
		queryStr        string
		queryStrBuilder strings.Builder

		jsonData *simplejson.Json
		err      error
	)

	if len(queryArgs) > 0 {
		if len(formDataStr) > 0 {
			queryStrBuilder = strings.Builder{}
			queryStrBuilder.Write(formDataStr)
			queryStrBuilder.WriteByte(38) //&
			queryStrBuilder.Write(queryArgs)
			queryStr = queryStrBuilder.String()
		} else {
			queryStr = strconv.B2S(queryArgs)
		}
	} else if len(formDataStr) > 0 {
		queryStr = strconv.B2S(formDataStr)
	}

	if len(body) > 0 {
		jsonData, err = simplejson.NewJson(body)
		if err == nil {
			params, err = jsonData.Map()
		}
		if err != nil {
			if len(queryStr) > 0 {
				params = map[string]any{}
				queryStrBuilder = strings.Builder{}
				queryStrBuilder.Write(body)
				queryStrBuilder.WriteByte(38) //&
				queryStrBuilder.WriteString(queryStr)
				queryStr = queryStrBuilder.String()
			} else {
				//queryArgs 是空的话就没直接返回
				return map[string]any{}
			}
		}
	} else if len(queryStr) == 0 {
		//两个参数都是空的话就直接返回
		return map[string]any{}
	}
	if params == nil {
		params = map[string]any{}
	}

	var (
		paramsSplit = strings.Split(queryStr, "&")
		index       int
		key         string
		value       string
	)
	for _, kv := range paramsSplit {
		if kv == "" {
			continue
		}
		index = strings.IndexAny(kv, "=")
		if index == -1 {
			continue
		}
		key = kv[:index]
		value = kv[(index + 1):]
		if value == "" {
			params[key] = ""
			continue
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			params[key] = value
			continue
		}
		jsonData, err = simplejson.NewJson(strconv.S2B(value))
		if err != nil {
			params[key] = value
			continue
		}

		paramValue := jsonData.Interface()
		if paramValue == nil {
			params[key] = value
			continue
		}
		switch reflect.TypeOf(paramValue).Kind() {
		case reflect.Map, reflect.Slice:
			params[key] = paramValue
		default:
			params[key] = value
		}
	}
	return
}

func (ctx *RequestCtx) SetInputs(data map[string]any) {
	(*fasthttp.RequestCtx)(ctx).SetUserValue("inputs", data)
}

func (ctx *RequestCtx) Inputs() map[string]any {
	inputs := (*fasthttp.RequestCtx)(ctx).UserValue("inputs")
	if inputs == nil {
		ctx.initParams()
		inputs = (*fasthttp.RequestCtx)(ctx).UserValue("inputs")
	}
	return inputs.(map[string]any)
}

func (ctx *RequestCtx) Input(key string, defaultVals ...any) any {
	input := ctx.Inputs()
	if input == nil {
		return nil
	}
	if value, ok := input[key]; ok {
		return value
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return nil
}

func (ctx *RequestCtx) InputsAssign(valPtr any) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("recover:" + fmt.Sprint(rec))
		}
	}()
	input := ctx.Inputs()
	if input == nil {
		return nil
	}
	marshal, err := json.Marshal(input)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(marshal, valPtr); err != nil {
		return err
	}
	return nil
}

func (ctx *RequestCtx) InputAssign(key string, valPtr any) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("recover:" + fmt.Sprint(rec))
		}
	}()
	input := ctx.Inputs()
	if input == nil {
		return nil
	}

	value, ok := input[key]
	if !ok {
		return nil
	}

	_reflect := reflect.TypeOf(valPtr)
	if _reflect.Kind() != reflect.Pointer {
		return errors.New("参数必须是指针类型")
	}
	if err = typeAssign(value, reflect.ValueOf(valPtr).Elem()); err != nil {
		return err
	}
	return nil
}

func typeAssign(value any, typeValue reflect.Value) (err error) {

	typeReflect := typeValue.Type()
	switch typeValue.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var _v int64
		switch raw := value.(type) {
		case json.Number:
			if _v, err = value.(json.Number).Int64(); err != nil {
				return err
			}
		case string:
			if _v, err = strconv2.ParseInt(raw, 10, 64); err != nil {
				return err
			}
		case int:
			_v = int64(raw)
		case int8:
			_v = int64(raw)
		case int16:
			_v = int64(raw)
		case int32:
			_v = int64(raw)
		case int64:
			_v = raw
		case uint:
			_v = int64(raw)
		case uint8:
			_v = int64(raw)
		case uint16:
			_v = int64(raw)
		case uint32:
			_v = int64(raw)
		case uint64:
			_v = int64(raw)
		case float32:
			_v = int64(raw)
		case float64:
			_v = int64(raw)
		case uintptr:
			_v = int64(raw)
		default:
			return NewDataError("数据类型错误", value)
		}
		_reflect := reflect.New(typeReflect).Elem()
		_reflect.SetInt(_v)
		typeValue.Set(_reflect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var _v uint64
		switch raw := value.(type) {
		case json.Number:
			if i, err := raw.Int64(); err != nil {
				return err
			} else {
				_v = uint64(i)
			}
		case string:
			if i, err := strconv2.ParseInt(raw, 10, 64); err != nil {
				return err
			} else {
				_v = uint64(i)
			}
		case int:
			_v = uint64(raw)
		case int8:
			_v = uint64(raw)
		case int16:
			_v = uint64(raw)
		case int32:
			_v = uint64(raw)
		case int64:
			_v = uint64(raw)
		case uint:
			_v = uint64(raw)
		case uint8:
			_v = uint64(raw)
		case uint16:
			_v = uint64(raw)
		case uint32:
			_v = uint64(raw)
		case uint64:
			_v = raw
		case float32:
			_v = uint64(raw)
		case float64:
			_v = uint64(raw)
		case uintptr:
			_v = uint64(raw)
		default:
			return NewDataError("数据类型错误", value)
		}
		_reflect := reflect.New(typeReflect).Elem()
		_reflect.SetUint(_v)
		typeValue.Set(_reflect)

	case reflect.Float64, reflect.Float32:
		var _v float64
		switch raw := value.(type) {
		case json.Number:
			if _v, err = raw.Float64(); err != nil {
				return err
			}
		case string:
			if _v, err = strconv2.ParseFloat(raw, 64); err != nil {
				return err
			}
		case int:
			_v = float64(raw)
		case int8:
			_v = float64(raw)
		case int16:
			_v = float64(raw)
		case int32:
			_v = float64(raw)
		case int64:
			_v = float64(raw)
		case uint:
			_v = float64(raw)
		case uint8:
			_v = float64(raw)
		case uint16:
			_v = float64(raw)
		case uint32:
			_v = float64(raw)
		case uint64:
			_v = float64(raw)
		case float32:
			_v = float64(raw)
		case float64:
			_v = raw
		case uintptr:
			_v = float64(raw)
		default:
			return NewDataError("数据类型错误", value)
		}
		_reflect := reflect.New(typeReflect).Elem()
		_reflect.SetFloat(_v)
		typeValue.Set(_reflect)

	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		marshal, err := helper.AnyToBytes(value)
		if err != nil {
			return err
		}
		_v := reflect.New(typeReflect).Interface()
		if err = json.Unmarshal(marshal, _v); err != nil {
			return err
		}
		typeValue.Set(reflect.ValueOf(_v).Elem())

	case reflect.String:
		_v, ok := value.(string)
		if !ok {
			_v, err = serialize.JsonEncode(value)
			if err != nil {
				return err
			}
		}
		typeValue.SetString(_v)
	case reflect.Bool:
		var b = false
		switch value.(type) {
		case string:
			_value := value.(string)
			if len(_value) != 0 && _value != "false" && _value != "0" {
				b = true
			}
		default:
			b = func() bool {
				defer func() {
					recover()
				}()
				return !reflect.ValueOf(value).IsZero()
			}()
		}
		typeValue.SetBool(b)

	default:
		return NewDataError("不受支持的数据类型", typeValue.Interface())
	}
	return nil
}

func (ctx *RequestCtx) Hijack(handler fasthttp.HijackHandler) {
	(*fasthttp.RequestCtx)(ctx).Hijack(handler)
}

func (ctx *RequestCtx) HijackSetNoResponse(noResponse bool) {
	(*fasthttp.RequestCtx)(ctx).HijackSetNoResponse(noResponse)
}

func (ctx *RequestCtx) Hijacked() bool {
	return (*fasthttp.RequestCtx)(ctx).Hijacked()
}

func (ctx *RequestCtx) SetUserValue(key interface{}, value interface{}) {
	(*fasthttp.RequestCtx)(ctx).SetUserValue(key, value)
}

func (ctx *RequestCtx) SetUserValueBytes(key []byte, value interface{}) {
	(*fasthttp.RequestCtx)(ctx).SetUserValueBytes(key, value)
}

func (ctx *RequestCtx) UserValue(key interface{}) interface{} {
	return (*fasthttp.RequestCtx)(ctx).UserValue(key)
}

func (ctx *RequestCtx) UserValueBytes(key []byte) interface{} {
	return (*fasthttp.RequestCtx)(ctx).UserValueBytes(key)
}

func (ctx *RequestCtx) VisitUserValues(visitor func([]byte, interface{})) {
	(*fasthttp.RequestCtx)(ctx).VisitUserValues(visitor)
}

func (ctx *RequestCtx) VisitUserValuesAll(visitor func(interface{}, interface{})) {
	(*fasthttp.RequestCtx)(ctx).VisitUserValuesAll(visitor)
}

func (ctx *RequestCtx) ResetUserValues() {
	(*fasthttp.RequestCtx)(ctx).ResetUserValues()
}

func (ctx *RequestCtx) RemoveUserValue(key interface{}) {
	(*fasthttp.RequestCtx)(ctx).RemoveUserValue(key)
}

func (ctx *RequestCtx) IsTLS() bool {
	return (*fasthttp.RequestCtx)(ctx).IsTLS()
}

func (ctx *RequestCtx) TLSConnectionState() *tls.ConnectionState {
	return (*fasthttp.RequestCtx)(ctx).TLSConnectionState()
}

func (ctx *RequestCtx) Conn() net.Conn {
	return (*fasthttp.RequestCtx)(ctx).Conn()
}

func (ctx *RequestCtx) String() string {
	return (*fasthttp.RequestCtx)(ctx).String()
}

func (ctx *RequestCtx) ID() uint64 {
	return (*fasthttp.RequestCtx)(ctx).ID()
}

func (ctx *RequestCtx) ConnID() uint64 {
	return (*fasthttp.RequestCtx)(ctx).ConnID()
}

func (ctx *RequestCtx) Time() time.Time {
	return (*fasthttp.RequestCtx)(ctx).Time()
}

func (ctx *RequestCtx) ConnTime() time.Time {
	return (*fasthttp.RequestCtx)(ctx).ConnTime()
}

func (ctx *RequestCtx) ConnRequestNum() uint64 {
	return (*fasthttp.RequestCtx)(ctx).ConnRequestNum()
}

func (ctx *RequestCtx) SetConnectionClose() {
	(*fasthttp.RequestCtx)(ctx).SetConnectionClose()
}

func (ctx *RequestCtx) SetStatusCode(statusCode int) {
	(*fasthttp.RequestCtx)(ctx).SetStatusCode(statusCode)
}

func (ctx *RequestCtx) SetContentType(contentType string) {
	(*fasthttp.RequestCtx)(ctx).SetContentType(contentType)
}

func (ctx *RequestCtx) SetContentTypeBytes(contentType []byte) {
	(*fasthttp.RequestCtx)(ctx).SetContentTypeBytes(contentType)
}

func (ctx *RequestCtx) RequestURI() []byte {
	return (*fasthttp.RequestCtx)(ctx).RequestURI()
}

func (ctx *RequestCtx) URI() *fasthttp.URI {
	return (*fasthttp.RequestCtx)(ctx).URI()
}

func (ctx *RequestCtx) Referer() []byte {
	return (*fasthttp.RequestCtx)(ctx).Referer()
}

func (ctx *RequestCtx) UserAgent() []byte {
	return (*fasthttp.RequestCtx)(ctx).UserAgent()
}

func (ctx *RequestCtx) Path() []byte {
	return (*fasthttp.RequestCtx)(ctx).Path()
}

func (ctx *RequestCtx) Host() []byte {
	return (*fasthttp.RequestCtx)(ctx).Host()
}

func (ctx *RequestCtx) QueryArgs() *fasthttp.Args {
	return (*fasthttp.RequestCtx)(ctx).QueryArgs()
}

func (ctx *RequestCtx) PostArgs() *fasthttp.Args {
	return (*fasthttp.RequestCtx)(ctx).PostArgs()
}

func (ctx *RequestCtx) MultipartForm() (*multipart.Form, error) {
	return (*fasthttp.RequestCtx)(ctx).MultipartForm()
}

func (ctx *RequestCtx) FormFile(key string) (*multipart.FileHeader, error) {
	return (*fasthttp.RequestCtx)(ctx).FormFile(key)
}

func (ctx *RequestCtx) FormValue(key string) []byte {
	return (*fasthttp.RequestCtx)(ctx).FormValue(key)
}

func (ctx *RequestCtx) IsGet() bool {
	return (*fasthttp.RequestCtx)(ctx).IsGet()
}

func (ctx *RequestCtx) IsPost() bool {
	return (*fasthttp.RequestCtx)(ctx).IsPost()
}

func (ctx *RequestCtx) IsPut() bool {
	return (*fasthttp.RequestCtx)(ctx).IsPut()
}

func (ctx *RequestCtx) IsDelete() bool {
	return (*fasthttp.RequestCtx)(ctx).IsDelete()
}

func (ctx *RequestCtx) IsConnect() bool {
	return (*fasthttp.RequestCtx)(ctx).IsConnect()
}

func (ctx *RequestCtx) IsOptions() bool {
	return (*fasthttp.RequestCtx)(ctx).IsOptions()
}

func (ctx *RequestCtx) IsTrace() bool {
	return (*fasthttp.RequestCtx)(ctx).IsTrace()
}

func (ctx *RequestCtx) IsPatch() bool {
	return (*fasthttp.RequestCtx)(ctx).IsPatch()
}

func (ctx *RequestCtx) Method() []byte {
	return (*fasthttp.RequestCtx)(ctx).Method()
}

func (ctx *RequestCtx) IsHead() bool {
	return (*fasthttp.RequestCtx)(ctx).IsHead()
}

func (ctx *RequestCtx) RemoteAddr() net.Addr {
	return (*fasthttp.RequestCtx)(ctx).RemoteAddr()
}

func (ctx *RequestCtx) SetRemoteAddr(remoteAddr net.Addr) {
	(*fasthttp.RequestCtx)(ctx).SetRemoteAddr(remoteAddr)
}

func (ctx *RequestCtx) LocalAddr() net.Addr {
	return (*fasthttp.RequestCtx)(ctx).LocalAddr()
}

func (ctx *RequestCtx) RemoteIP() net.IP {
	return (*fasthttp.RequestCtx)(ctx).RemoteIP()
}

func (ctx *RequestCtx) LocalIP() net.IP {
	return (*fasthttp.RequestCtx)(ctx).LocalIP()
}

func (ctx *RequestCtx) Error(msg string, statusCode int) {
	(*fasthttp.RequestCtx)(ctx).Error(msg, statusCode)
}

func (ctx *RequestCtx) Success(contentType string, body []byte) {
	(*fasthttp.RequestCtx)(ctx).Success(contentType, body)
}

func (ctx *RequestCtx) SuccessString(contentType, body string) {
	(*fasthttp.RequestCtx)(ctx).SuccessString(contentType, body)
}

func (ctx *RequestCtx) Redirect(uri string, statusCode int) {
	(*fasthttp.RequestCtx)(ctx).Redirect(uri, statusCode)
}

func (ctx *RequestCtx) RedirectBytes(uri []byte, statusCode int) {
	(*fasthttp.RequestCtx)(ctx).RedirectBytes(uri, statusCode)
}

func (ctx *RequestCtx) SetBody(body []byte) {
	(*fasthttp.RequestCtx)(ctx).SetBody(body)
}

func (ctx *RequestCtx) SetBodyString(body string) {
	(*fasthttp.RequestCtx)(ctx).SetBodyString(body)
}

func (ctx *RequestCtx) ResetBody() {
	(*fasthttp.RequestCtx)(ctx).ResetBody()
}

func (ctx *RequestCtx) SendFile(path string) {
	(*fasthttp.RequestCtx)(ctx).SendFile(path)
}

func (ctx *RequestCtx) SendFileBytes(path []byte) {
	(*fasthttp.RequestCtx)(ctx).SendFileBytes(path)
}

func (ctx *RequestCtx) IfModifiedSince(lastModified time.Time) bool {
	return (*fasthttp.RequestCtx)(ctx).IfModifiedSince(lastModified)
}

func (ctx *RequestCtx) NotModified() {
	(*fasthttp.RequestCtx)(ctx).NotModified()
}

func (ctx *RequestCtx) NotFound() {
	(*fasthttp.RequestCtx)(ctx).NotFound()
}

func (ctx *RequestCtx) Write(p []byte) (int, error) {
	return (*fasthttp.RequestCtx)(ctx).Write(p)
}

func (ctx *RequestCtx) WriteString(s string) (int, error) {
	return (*fasthttp.RequestCtx)(ctx).WriteString(s)
}

func (ctx *RequestCtx) PostBody() []byte {
	return (*fasthttp.RequestCtx)(ctx).PostBody()
}

func (ctx *RequestCtx) SetBodyStream(bodyStream io.Reader, bodySize int) {
	(*fasthttp.RequestCtx)(ctx).SetBodyStream(bodyStream, bodySize)
}

func (ctx *RequestCtx) SetBodyStreamWriter(sw fasthttp.StreamWriter) {
	(*fasthttp.RequestCtx)(ctx).SetBodyStreamWriter(sw)
}

func (ctx *RequestCtx) IsBodyStream() bool {
	return (*fasthttp.RequestCtx)(ctx).IsBodyStream()
}

func (ctx *RequestCtx) Logger() fasthttp.Logger {
	return (*fasthttp.RequestCtx)(ctx).Logger()
}

func (ctx *RequestCtx) TimeoutError(msg string) {
	(*fasthttp.RequestCtx)(ctx).TimeoutError(msg)
}

func (ctx *RequestCtx) TimeoutErrorWithCode(msg string, statusCode int) {
	(*fasthttp.RequestCtx)(ctx).TimeoutErrorWithCode(msg, statusCode)
}

func (ctx *RequestCtx) TimeoutErrorWithResponse(resp *fasthttp.Response) {
	(*fasthttp.RequestCtx)(ctx).TimeoutErrorWithResponse(resp)
}

func (ctx *RequestCtx) LastTimeoutErrorResponse() *fasthttp.Response {
	return (*fasthttp.RequestCtx)(ctx).LastTimeoutErrorResponse()
}

func (ctx *RequestCtx) Init2(conn net.Conn, logger fasthttp.Logger, reduceMemoryUsage bool) {
	(*fasthttp.RequestCtx)(ctx).Init2(conn, logger, reduceMemoryUsage)
}

func (ctx *RequestCtx) Init(req *fasthttp.Request, remoteAddr net.Addr, logger fasthttp.Logger) {
	(*fasthttp.RequestCtx)(ctx).Init(req, remoteAddr, logger)
}

func (ctx *RequestCtx) Deadline() (deadline time.Time, ok bool) {
	return (*fasthttp.RequestCtx)(ctx).Deadline()
}

func (ctx *RequestCtx) Done() <-chan struct{} {
	return (*fasthttp.RequestCtx)(ctx).Done()
}

func (ctx *RequestCtx) Err() error {
	return (*fasthttp.RequestCtx)(ctx).Err()
}

func (ctx *RequestCtx) Value(key interface{}) interface{} {
	return (*fasthttp.RequestCtx)(ctx).Value(key)
}
