package types

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fushiliang321/go-core/helper"
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

const gzipMinSize = 10000 //触发gzip压缩的最小长度

type RequestCtx fasthttp.RequestCtx

// any转bytes
func anyToBytes(data any) (bts []byte, err error) {
	switch data.(type) {
	case string:
		bts = strconv.S2B(data.(string))
	case *string:
		bts = strconv.S2B(*(data.(*string)))
	case []byte:
		return data.([]byte), nil
	case *[]byte:
		bts = *data.(*[]byte)
	case byte:
		bts = []byte{data.(byte)}
	case *byte:
		bts = []byte{*(data.(*byte))}
	default:
		bts, err = json.Marshal(data)
	}
	return
}

func (ctx *RequestCtx) Raw() *fasthttp.RequestCtx {
	return (*fasthttp.RequestCtx)(ctx)
}

func (ctx *RequestCtx) WriteAny(data any) (int, error) {
	bytes, err := anyToBytes(data)
	if err != nil {
		log.Printf("write data err:%s\n", err)
		return 0, err
	}
	if len(bytes) > gzipMinSize {
		ctx.Response.Header.Add("Content-Encoding", "gzip")
		return ctx.Raw().Write(fasthttp.AppendGzipBytes([]byte{}, bytes))
	}
	return ctx.Raw().Write(bytes)
}

func (ctx *RequestCtx) initParams() {
	var (
		body        = ctx.Raw().PostBody()
		paramsStr   = ctx.Raw().QueryArgs().QueryString()
		formDataStr = ctx.Raw().PostArgs().QueryString()
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
	ctx.Raw().SetUserValue("inputs", data)
}

func (ctx *RequestCtx) Inputs() map[string]any {
	inputs := ctx.Raw().UserValue("inputs")
	if inputs == nil {
		ctx.initParams()
		inputs = ctx.Raw().UserValue("inputs")
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
		fmt.Println(1)
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
		fmt.Println(_reflect.Interface())
		typeValue.Set(_reflect)
		fmt.Println(typeValue.Interface())

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
			_v, err = helper.JsonEncode(value)
			if err != nil {
				return err
			}
		}
		typeValue.SetString(_v)

	default:
		return NewDataError("不受支持的数据类型", typeValue.Interface())
	}
	return nil
}

func (ctx *RequestCtx) Hijack(handler fasthttp.HijackHandler) {
	ctx.Raw().Hijack(handler)
}

func (ctx *RequestCtx) HijackSetNoResponse(noResponse bool) {
	ctx.Raw().HijackSetNoResponse(noResponse)
}

func (ctx *RequestCtx) Hijacked() bool {
	return ctx.Raw().Hijacked()
}

func (ctx *RequestCtx) SetUserValue(key interface{}, value interface{}) {
	ctx.Raw().SetUserValue(key, value)
}

func (ctx *RequestCtx) SetUserValueBytes(key []byte, value interface{}) {
	ctx.Raw().SetUserValueBytes(key, value)
}

func (ctx *RequestCtx) UserValue(key interface{}) interface{} {
	return ctx.Raw().UserValue(key)
}
func (ctx *RequestCtx) UserValueBytes(key []byte) interface{} {
	return ctx.Raw().UserValueBytes(key)
}
func (ctx *RequestCtx) VisitUserValues(visitor func([]byte, interface{})) {
	ctx.Raw().VisitUserValues(visitor)
}

func (ctx *RequestCtx) VisitUserValuesAll(visitor func(interface{}, interface{})) {
	ctx.Raw().VisitUserValuesAll(visitor)
}

func (ctx *RequestCtx) ResetUserValues() {
	ctx.Raw().ResetUserValues()
}

func (ctx *RequestCtx) RemoveUserValue(key interface{}) {
	ctx.Raw().RemoveUserValue(key)
}

func (ctx *RequestCtx) IsTLS() bool {
	return ctx.Raw().IsTLS()
}

func (ctx *RequestCtx) TLSConnectionState() *tls.ConnectionState {
	return ctx.Raw().TLSConnectionState()
}

func (ctx *RequestCtx) Conn() net.Conn {
	return ctx.Raw().Conn()
}
func (ctx *RequestCtx) String() string {
	return ctx.Raw().String()
}
func (ctx *RequestCtx) ID() uint64 {
	return ctx.Raw().ID()
}
func (ctx *RequestCtx) ConnID() uint64 {
	return ctx.Raw().ConnID()
}
func (ctx *RequestCtx) Time() time.Time {
	return ctx.Raw().Time()
}
func (ctx *RequestCtx) ConnTime() time.Time {
	return ctx.Raw().ConnTime()
}
func (ctx *RequestCtx) ConnRequestNum() uint64 {
	return ctx.Raw().ConnRequestNum()
}
func (ctx *RequestCtx) SetConnectionClose() {
	ctx.Raw().SetConnectionClose()
}
func (ctx *RequestCtx) SetStatusCode(statusCode int) {
	ctx.Raw().SetStatusCode(statusCode)
}
func (ctx *RequestCtx) SetContentType(contentType string) {
	ctx.Raw().SetContentType(contentType)
}
func (ctx *RequestCtx) SetContentTypeBytes(contentType []byte) {
	ctx.Raw().SetContentTypeBytes(contentType)
}
func (ctx *RequestCtx) RequestURI() []byte {
	return ctx.Raw().RequestURI()
}
func (ctx *RequestCtx) URI() *fasthttp.URI {
	return ctx.Raw().URI()
}

func (ctx *RequestCtx) Referer() []byte {
	return ctx.Raw().Referer()
}

func (ctx *RequestCtx) UserAgent() []byte {
	return ctx.Raw().UserAgent()
}

func (ctx *RequestCtx) Path() []byte {
	return ctx.Raw().Path()
}

func (ctx *RequestCtx) Host() []byte {
	return ctx.Raw().Host()
}

func (ctx *RequestCtx) QueryArgs() *fasthttp.Args {
	return ctx.Raw().QueryArgs()
}

func (ctx *RequestCtx) PostArgs() *fasthttp.Args {
	return ctx.Raw().PostArgs()
}

func (ctx *RequestCtx) MultipartForm() (*multipart.Form, error) {
	return ctx.Raw().MultipartForm()
}

func (ctx *RequestCtx) FormFile(key string) (*multipart.FileHeader, error) {
	return ctx.Raw().FormFile(key)
}

func (ctx *RequestCtx) FormValue(key string) []byte {
	return ctx.Raw().FormValue(key)
}

func (ctx *RequestCtx) IsGet() bool {
	return ctx.Raw().IsGet()
}

func (ctx *RequestCtx) IsPost() bool {
	return ctx.Raw().IsPost()
}
func (ctx *RequestCtx) IsPut() bool {
	return ctx.Raw().IsPut()
}
func (ctx *RequestCtx) IsDelete() bool {
	return ctx.Raw().IsDelete()
}
func (ctx *RequestCtx) IsConnect() bool {
	return ctx.Raw().IsConnect()
}
func (ctx *RequestCtx) IsOptions() bool {
	return ctx.Raw().IsOptions()
}
func (ctx *RequestCtx) IsTrace() bool {
	return ctx.Raw().IsTrace()
}
func (ctx *RequestCtx) IsPatch() bool {
	return ctx.Raw().IsPatch()
}
func (ctx *RequestCtx) Method() []byte {
	return ctx.Raw().Method()
}
func (ctx *RequestCtx) IsHead() bool {
	return ctx.Raw().IsHead()
}
func (ctx *RequestCtx) RemoteAddr() net.Addr {
	return ctx.Raw().RemoteAddr()
}
func (ctx *RequestCtx) SetRemoteAddr(remoteAddr net.Addr) {
	ctx.Raw().SetRemoteAddr(remoteAddr)
}
func (ctx *RequestCtx) LocalAddr() net.Addr {
	return ctx.Raw().LocalAddr()
}
func (ctx *RequestCtx) RemoteIP() net.IP {
	return ctx.Raw().RemoteIP()
}
func (ctx *RequestCtx) LocalIP() net.IP {
	return ctx.Raw().LocalIP()
}
func (ctx *RequestCtx) Error(msg string, statusCode int) {
	ctx.Raw().Error(msg, statusCode)
}
func (ctx *RequestCtx) Success(contentType string, body []byte) {
	ctx.Raw().Success(contentType, body)
}
func (ctx *RequestCtx) SuccessString(contentType, body string) {
	ctx.Raw().SuccessString(contentType, body)
}
func (ctx *RequestCtx) Redirect(uri string, statusCode int) {
	ctx.Raw().Redirect(uri, statusCode)
}
func (ctx *RequestCtx) RedirectBytes(uri []byte, statusCode int) {
	ctx.Raw().RedirectBytes(uri, statusCode)
}
func (ctx *RequestCtx) SetBody(body []byte) {
	ctx.Raw().SetBody(body)
}
func (ctx *RequestCtx) SetBodyString(body string) {
	ctx.Raw().SetBodyString(body)
}
func (ctx *RequestCtx) ResetBody() {
	ctx.Raw().ResetBody()
}
func (ctx *RequestCtx) SendFile(path string) {
	ctx.Raw().SendFile(path)
}

func (ctx *RequestCtx) SendFileBytes(path []byte) {
	ctx.Raw().SendFileBytes(path)
}

func (ctx *RequestCtx) IfModifiedSince(lastModified time.Time) bool {
	return ctx.Raw().IfModifiedSince(lastModified)
}

func (ctx *RequestCtx) NotModified() {
	ctx.Raw().NotModified()
}

func (ctx *RequestCtx) NotFound() {
	ctx.Raw().NotFound()
}

func (ctx *RequestCtx) Write(p []byte) (int, error) {
	return ctx.Raw().Write(p)
}

func (ctx *RequestCtx) WriteString(s string) (int, error) {
	return ctx.Raw().WriteString(s)
}

func (ctx *RequestCtx) PostBody() []byte {
	return ctx.Raw().PostBody()
}

func (ctx *RequestCtx) SetBodyStream(bodyStream io.Reader, bodySize int) {
	ctx.Raw().SetBodyStream(bodyStream, bodySize)
}

func (ctx *RequestCtx) SetBodyStreamWriter(sw fasthttp.StreamWriter) {
	ctx.Raw().SetBodyStreamWriter(sw)
}

func (ctx *RequestCtx) IsBodyStream() bool {
	return ctx.Raw().IsBodyStream()
}

func (ctx *RequestCtx) Logger() fasthttp.Logger {
	return ctx.Raw().Logger()
}

func (ctx *RequestCtx) TimeoutError(msg string) {
	ctx.Raw().TimeoutError(msg)
}

func (ctx *RequestCtx) TimeoutErrorWithCode(msg string, statusCode int) {
	ctx.Raw().TimeoutErrorWithCode(msg, statusCode)
}

func (ctx *RequestCtx) TimeoutErrorWithResponse(resp *fasthttp.Response) {
	ctx.Raw().TimeoutErrorWithResponse(resp)
}

func (ctx *RequestCtx) LastTimeoutErrorResponse() *fasthttp.Response {
	return ctx.Raw().LastTimeoutErrorResponse()
}

func (ctx *RequestCtx) Init2(conn net.Conn, logger fasthttp.Logger, reduceMemoryUsage bool) {
	ctx.Raw().Init2(conn, logger, reduceMemoryUsage)
}

func (ctx *RequestCtx) Init(req *fasthttp.Request, remoteAddr net.Addr, logger fasthttp.Logger) {
	ctx.Raw().Init(req, remoteAddr, logger)
}

func (ctx *RequestCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.Raw().Deadline()
}

func (ctx *RequestCtx) Done() <-chan struct{} {
	return ctx.Raw().Done()
}

func (ctx *RequestCtx) Err() error {
	return ctx.Raw().Err()
}

func (ctx *RequestCtx) Value(key interface{}) interface{} {
	return ctx.Raw().Value(key)
}
