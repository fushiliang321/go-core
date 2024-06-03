package string

import (
	"strings"
)

const (
	underlineByte = byte('_')
	diffValue     = uint8('a' - 'A')
)

// 转为蛇形字符串
func SnakeString(s string) string {
	var (
		num = len(s)
		d   uint8
		b   = strings.Builder{}
	)
	for i := 0; i < num; i++ {
		d = s[i]
		if d >= 'A' && d <= 'Z' {
			if i > 0 {
				b.WriteByte(underlineByte)
			}
			b.WriteByte(d + diffValue)
		} else {
			b.WriteByte(d)
		}

	}
	return b.String()
}
