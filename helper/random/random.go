package random

import (
	"math/rand"
	"time"
	"unsafe"
)

type letterType struct {
	bytes string
	bits  int64
	mask  int64
	max   int64
}

const (
	RangeLetter0 = "0"
	RangeLetter1 = "1"

	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterBytes1 = "abcdefghijklmnopqrstuvwxyz1234567890"

	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = letterIdxMask / letterIdxBits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 取随机数
func RangeRand(min, max int) int {
	if min > max {
		panic("the min is greater than max!")
	}
	dif := max - min
	return min + rand.Intn(dif)
}

// 获取指定长度的随机字符串
func randString(n int, bytes string) string {
	bytesLen := len(bytes)
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < bytesLen {
			b[i] = bytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// 获取指定长度的随机字符串,可自定义字符串
func RandString(n int, bytes ...string) string {
	switch len(bytes) {
	case 0:
		return randString(n, letterBytes)
	case 1:
		switch bytes[0] {
		case RangeLetter0:
			return randString(n, letterBytes)
		case RangeLetter1:
			return randString(n, letterBytes1)
		}
	}
	var (
		letterIdxBits = 3
		letterIdxMask = 1<<letterIdxBits - 1
		letterStr     = bytes[0]
	)
	bytesLen := len(letterStr)

	for ; letterIdxMask < bytesLen; letterIdxBits++ {
		letterIdxMask = 1<<letterIdxBits - 1
	}
	letterIdxMax := letterIdxMask / letterIdxBits
	letterIdxMaskInt64 := int64(letterIdxMask)
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMaskInt64); idx < bytesLen {
			b[i] = letterStr[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
