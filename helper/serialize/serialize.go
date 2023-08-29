package serialize

import (
	"bytes"
	"encoding/json"
	"github.com/savsgio/gotils/strconv"
)

// json字符串解码
func JsonDecode(str string, v any) error {
	d := json.NewDecoder(bytes.NewReader(strconv.S2B(str)))
	d.UseNumber()
	return d.Decode(&v)
}

// json编码
func JsonEncode(v any) (string, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return strconv.B2S(marshal), nil
}
