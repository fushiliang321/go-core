package server

import (
	"github.com/fushiliang321/go-core/jsonRpcHttp/context"
	"github.com/fushiliang321/jsonrpc/common"
	"reflect"
)

type Server struct {
	common.Server
}

func (svr *Server) Handler(b []byte) any {
	data, err := common.ParseRequestBody(b)
	if err != nil {
		return jsonE(nil, common.JsonRpc, common.ParseError)
	}
	var res any
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		var resList []any
		for _, v := range data.([]any) {
			r := svr.SingleHandler(v.(map[string]any))
			resList = append(resList, r)
		}
		res = resList
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		r := svr.SingleHandler(data.(map[string]any))
		res = r
	} else {
		return jsonE(nil, common.JsonRpc, common.InvalidRequest)
	}
	return res
}

func (svr *Server) SingleHandler(jsonMap map[string]any) any {
	if ctx, ok := jsonMap["context"]; ok {
		if ctx, ok := ctx.(map[string]any); ok {
			context.SetBatch(ctx)
		}
	}
	return svr.Server.SingleHandler(jsonMap)
}

func jsonE(id any, jsonRpc string, errCode int) any {
	return common.E(id, jsonRpc, errCode)
}
