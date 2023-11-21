package redis

// 设置
func HSet(key string, values ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.HSet(_ctx, key, values...).Result()
}

// 批量设置
func HMSet(key string, values ...any) (bool, error) {
	c, err := client()
	if err != nil {
		return false, err
	}
	return c.HMSet(_ctx, key, values...).Result()
}

// 获取某个元素
func HGet(key, field string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.HGet(_ctx, key, field).Result()
}

// 获取某个元素
func HGetAll(key string) (map[string]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.HGetAll(_ctx, key).Result()
}

// 删除某个元素
func HDel(key string, fields ...string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.HDel(_ctx, key, fields...).Result()
}

// 判断元素是否存在
func HExists(key, field string) (bool, error) {
	c, err := client()
	if err != nil {
		return false, err
	}
	return c.HExists(_ctx, key, field).Result()
}

// 获取长度
func HLen(key string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.HLen(_ctx, key).Result()
}
