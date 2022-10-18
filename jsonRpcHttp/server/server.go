package server

import (
	"encoding/json"
	"github.com/fushiliang321/go-core/jsonRpcHttp/context"
	"github.com/iloveswift/go-jsonrpc/common"
	"reflect"
)

type Server struct {
	common.Server
}

func (svr *Server) Handler(b []byte) []byte {
	data, err := common.ParseRequestBody(b)
	if err != nil {
		return jsonE(nil, common.JsonRpc, common.ParseError)
	}
	var res interface{}
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		var resList []interface{}
		for _, v := range data.([]interface{}) {
			r := svr.SingleHandler(v.(map[string]interface{}))
			resList = append(resList, r)
		}
		res = resList
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		r := svr.SingleHandler(data.(map[string]interface{}))
		res = r
	} else {
		return jsonE(nil, common.JsonRpc, common.InvalidRequest)
	}

	response, _ := json.Marshal(res)
	return response
}

func (svr *Server) SingleHandler(jsonMap map[string]interface{}) interface{} {
	if ctx, ok := jsonMap["context"]; ok {
		if ctx, ok := ctx.(map[string]any); ok {
			context.SetBatch(ctx)
		}
	}
	return svr.Server.SingleHandler(jsonMap)
}

func jsonE(id interface{}, jsonRpc string, errCode int) []byte {
	e, _ := json.Marshal(common.E(id, jsonRpc, errCode))
	return e
}
