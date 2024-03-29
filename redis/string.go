package redis

import (
	"encoding/json"
	"errors"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/savsgio/gotils/strconv"
	"reflect"
	"time"
)

// 设置指定 key 的值
func Set(key string, val any, expiration ...int64) error {
	var (
		exp    time.Duration
		valStr string
	)
	if len(expiration) > 0 && expiration[0] > 0 {
		exp = time.Duration(expiration[0]) * time.Second
	}
	switch _val := val.(type) {
	case string:
		valStr = _val
	case *string:
		valStr = *_val
	case byte:
		valStr = string(_val)
	case *byte:
		valStr = string(*_val)
	case []byte:
		valStr = strconv.B2S(_val)
	case *[]byte:
		valStr = strconv.B2S(*_val)
	default:
		valBytes, err := json.Marshal(val)
		if err != nil {
			logger.Warn("redis set serialization error:", err)
			return err
		}
		valStr = strconv.B2S(valBytes)
	}
	c, err := client()
	if err != nil {
		return err
	}
	err = c.Set(_ctx, key, valStr, exp).Err()
	if err != nil {
		logger.Warn("redis set error:", err.Error())
		return err
	}
	return nil
}

// 获取指定 key 的字符串值
func GetString(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	val, err := c.Get(_ctx, key).Result()
	if err != nil {
		if err != Nil {
			//没取到值的时候不打印错误信息
			logger.Warn("redis get error:", err.Error())
		}
		return "", err
	}
	return val, nil
}

// 获取指定 key 的指定类型值
func Get[t any](key string) (*t, error) {
	var (
		res  t
		kind = reflect.TypeOf(res).Kind()
	)
	if kind == reflect.Ptr {
		logger.Warn("redis get error：", "the type cannot be a pointer")
		return nil, errors.New("the type cannot be a pointer")
	}
	v, err := GetString(key)
	if err != nil {
		return nil, err
	}
	if kind == reflect.String {
		res = any(v).(t)
		return &res, nil
	}
	err = json.Unmarshal(strconv.S2B(v), &res)
	if err != nil {
		logger.Warn("redis get deserialization error：", err.Error())
		return nil, err
	}
	return &res, nil
}

// 自增
func Inc(key string, v ...int64) (result int64, err error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	if len(v) > 0 {
		return c.IncrBy(_ctx, key, v[0]).Result()
	}
	return c.Incr(_ctx, key).Result()
}

// 自减
func Dec(key string, v ...int64) (result int64, err error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	if len(v) > 0 {
		return c.DecrBy(_ctx, key, v[0]).Result()
	}
	return c.Decr(_ctx, key).Result()
}

// 字符串截取
func GetRange(key string, start, end int64) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.GetRange(_ctx, key, start, end).Result()
}

// 追加到值的末尾
func Append(key, value string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.Append(_ctx, key, value).Result()
}

// 返回 key 所储存的字符串值的长度
func StrLen(key string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.StrLen(_ctx, key).Result()
}
