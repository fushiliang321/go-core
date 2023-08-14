package types

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"reflect"
	"strings"
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
	}
	return inputs.(map[string]any)
}

func (ctx *RequestCtx) Input(key string, defaultVal ...any) any {
	input := ctx.Inputs()
	if input == nil {
		return nil
	}
	if value, ok := input[key]; ok {
		return value
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return nil
}
