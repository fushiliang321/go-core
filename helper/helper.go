package helper

import (
	"bytes"
	"encoding/json"
	"github.com/savsgio/gotils/strconv"
	"reflect"
)

// map转struct
func MapToStruct[_type interface {
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

// 获取结构体字段
func GetStructFields(v any) []string {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Struct {
		typ := value.Type()
		numField := typ.NumField()
		fields := make([]string, numField)
		for i := 0; i < numField; i++ {
			fields[i] = typ.Field(i).Name
		}
		return fields
	}
	return nil
}

// any转bytes
func AnyToBytes(data any) ([]byte, error) {
	switch data := data.(type) {
	case string:
		return strconv.S2B(data), nil
	case *string:
		if data == nil {
			return nil, nil
		}
		return strconv.S2B(*data), nil
	case []byte:
		return data, nil
	case *[]byte:
		if data == nil {
			return nil, nil
		}
		return *data, nil
	case byte:
		return []byte{data}, nil
	case *byte:
		if data == nil {
			return nil, nil
		}
		return []byte{*data}, nil
	default:
		return json.Marshal(data)
	}
}

// any转string
func AnyToString(data any) (string, error) {
	switch data := data.(type) {
	case string:
		return data, nil
	case *string:
		if data == nil {
			return "", nil
		}
		return *data, nil
	case []byte:
		return strconv.B2S(data), nil
	case *[]byte:
		if data == nil {
			return "", nil
		}
		return strconv.B2S(*data), nil
	case byte:
		return strconv.B2S([]byte{data}), nil
	case *byte:
		if data == nil {
			return "", nil
		}
		return strconv.B2S([]byte{*data}), nil
	default:
		marshal, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return strconv.B2S(marshal), nil
	}
}
