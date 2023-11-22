package redis

import (
	"errors"
	"time"
)

// 根据正则获取keys
func Keys(pattern string) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.Keys(_ctx, pattern).Result()
}

// 获取key对应值得类型
func Type(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.Type(_ctx, key).Result()
}

// 删除缓存项
func Del(keys ...string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.Del(_ctx, keys...).Result()
}

// 检测缓存项是否存在
func Exists(keys ...string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.Exists(_ctx, keys...).Result()
}

// 设置过期时间，以秒计
func Expire(key string, exp int64) (bool, error) {
	c, err := client()
	if err != nil {
		return false, err
	}
	if exp < 1 {
		return false, errors.New("exp 值无效")
	}
	return c.Expire(_ctx, key, time.Duration(exp)*time.Second).Result()
}

// 设置过期时间，指定时间点
func ExpireAt(key string, tm time.Time) (bool, error) {
	c, err := client()
	if err != nil {
		return false, err
	}
	return c.ExpireAt(_ctx, key, tm).Result()
}

// 获取某个键的剩余有效期，以秒为单位
func TTL(key string) (time.Duration, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.TTL(_ctx, key).Result()
}

// 获取某个键的剩余有效期，以毫秒为单位
func PTTL(key string) (time.Duration, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.PTTL(_ctx, key).Result()
}
