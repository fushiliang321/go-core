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
	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		var resList []any
		for _, v := range data.([]any) {
			r := svr.SingleHandler(v.(map[string]any))
			resList = append(resList, r)
		}
		return resList
	case reflect.Map:
		return svr.SingleHandler(data.(map[string]any))
	default:
		return jsonE(nil, common.JsonRpc, common.InvalidRequest)
	}
}

func (svr *Server) SingleHandler(jsonMap map[string]any) any {
	if ctx, ok := jsonMap["context"]; ok {
		if ctx, ok := ctx.(map[string]any); ok && ctx != nil {
			context.SetBatch(ctx)
		}
	}
	return svr.Server.SingleHandler(jsonMap)
}

func jsonE(id any, jsonRpc string, errCode int) any {
	return common.E(id, jsonRpc, errCode)
}
