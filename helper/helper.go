package helper

import (
	"bytes"
	"encoding/json"
	"github.com/savsgio/gotils/strconv"
	"reflect"
)

var cacheData = map[string]any{}

// map转struc
func MapToStruc[_type interface {
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
func GetStructFields(v any) (fields []string) {
	if reflect.ValueOf(v).Kind() == reflect.Struct {
		reflectType := reflect.TypeOf(v)
		numField := reflectType.NumField()
		for i := 0; i < numField; i++ {
			fields = append(fields, reflectType.Field(i).Name)
		}
	}
	return
}

// any转bytes
func AnyToBytes(data any) (*[]byte, error) {
	var bts []byte
	var err error
	switch data.(type) {
	case string:
		bts = strconv.S2B(data.(string))
	case *string:
		bts = strconv.S2B(*(data.(*string)))
	case []byte:
		bts = data.([]byte)
	case *[]byte:
		return data.(*[]byte), nil
	case byte:
		bts = []byte{data.(byte)}
	case *byte:
		bts = []byte{*(data.(*byte))}
	default:
		bts, err = json.Marshal(data)
	}
	return &bts, err
}
