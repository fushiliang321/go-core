package helper

import (
	"github.com/fushiliang321/go-core/router/types"
	"reflect"
)

func Input[T any](ctx *types.RequestCtx, key string, defaultVal ...T) (*T, error) {
	var v T
	_reflect := reflect.New(reflect.TypeOf(v)).Elem()
	if len(defaultVal) > 0 {
		_reflect.Set(reflect.ValueOf(defaultVal[0]))
	}
	v = _reflect.Interface().(T)
	if err := ctx.InputAssign(key, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
