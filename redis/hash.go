package redis

// 设置
func HSet(key string, values ...any) (int64, error) {
	return client().HSet(_ctx, key, values...).Result()
}

// 批量设置
func HMSet(key string, values ...any) (bool, error) {
	return client().HMSet(_ctx, key, values...).Result()
}

// 获取某个元素
func HGet(key, field string) (string, error) {
	return client().HGet(_ctx, key, field).Result()
}

// 获取某个元素
func HGetAll(key string) (map[string]string, error) {
	return client().HGetAll(_ctx, key).Result()
}

// 删除某个元素
func HDel(key string, fields ...string) (int64, error) {
	return client().HDel(_ctx, key, fields...).Result()
}

// 判断元素是否存在
func HExists(key, field string) (bool, error) {
	return client().HExists(_ctx, key, field).Result()
}

// 获取长度
func HLen(key string) (int64, error) {
	return client().HLen(_ctx, key).Result()
}
