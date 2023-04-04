package helper

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/savsgio/gotils/strconv"
)

// 字符串md5
func MD5(v string) string {
	d := strconv.S2B(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
