package helper

import (
	"encoding/json"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/router/types"
	"reflect"
	strconv2 "strconv"
	"unsafe"
)

type DataError struct {
	text string
	data any
}

func (err *DataError) Error() string {
	return err.text
}
func (err *DataError) Data() any {
	return err.data
}

func NewDataError(text string, data any) *DataError {
	return &DataError{
		text: text,
		data: data,
	}
}

func Input[T any](ctx *types.RequestCtx, key string, defaultVal ...any) (*T, error) {
	var (
		v   T
		err error
	)
	value := ctx.Input(key, defaultVal[0])
	if value == nil {
		return nil, nil
	}
	_reflect := reflect.ValueOf(v)
	switch _reflect.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var _v int64
		switch raw := value.(type) {
		case json.Number:
			if _v, err = value.(json.Number).Int64(); err != nil {
				return nil, err
			}
		case string:
			if _v, err = strconv2.ParseInt(raw, 10, 64); err != nil {
				return nil, err
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
			return nil, NewDataError("数据类型错误", value)
		}
		_reflect.SetInt(_v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var _v uint64
		switch raw := value.(type) {
		case json.Number:
			if i, err := raw.Int64(); err != nil {
				return nil, err
			} else {
				_v = uint64(i)
			}
		case string:
			if i, err := strconv2.ParseInt(raw, 10, 64); err != nil {
				return nil, err
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
			return nil, NewDataError("数据类型错误", value)
		}
		_reflect.SetUint(_v)

	case reflect.Float64, reflect.Float32:
		var _v float64
		switch raw := value.(type) {
		case json.Number:
			if _v, err = raw.Float64(); err != nil {
				return nil, err
			}
		case string:
			if _v, err = strconv2.ParseFloat(raw, 64); err != nil {
				return nil, err
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
			return nil, NewDataError("数据类型错误", value)
		}
		_reflect.SetFloat(_v)

	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		if reflect.TypeOf(value).Kind() == _reflect.Kind() {
			_reflect.SetPointer(unsafe.Pointer(&value))
			break
		}
		marshal, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(marshal, &v); err != nil {
			return nil, err
		}

	case reflect.String:
		if _v, ok := value.(string); ok {
			_reflect.SetString(_v)
		} else {
			_v, err = helper.JsonEncode(value)
			if err != nil {
				return nil, err
			}
			_reflect.SetString(_v)
		}

	default:
		return nil, nil
	}
	return &v, nil
}
