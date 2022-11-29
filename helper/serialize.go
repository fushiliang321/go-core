package helper

import (
	"bytes"
	"encoding/json"
)

// json字符串解码
func JsonDecode(str string, v any) error {
	d := json.NewDecoder(bytes.NewReader([]byte(str)))
	d.UseNumber()
	return d.Decode(&v)
}

// json编码
func JsonEncode(v any) (string, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
