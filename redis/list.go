package redis

// 将元素从左侧压入链表
func LPush(key string, values ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.LPush(_ctx, key, values...).Result()
}

// 将元素从右侧压入链表
func RPush(key string, values ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.RPush(_ctx, key, values...).Result()
}

// 在某个位置插入新元素
func LInsert(key, op string, pivot, value any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.LInsert(_ctx, key, op, pivot, value).Result()
}

// 设置某个元素的值
func LSet(key string, index int64, value any) error {
	c, err := client()
	if err != nil {
		return err
	}
	return c.LSet(_ctx, key, index, value).Err()
}

// 获取链表元素个数
func LLen(key string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.LLen(_ctx, key).Result()
}

// 获取链表下标对应的元素
func LIndex(key string, index int64) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.LIndex(_ctx, key, index).Result()
}

// 获取某个选定范围的元素集
func LRange(key string, start, stop int64) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.LRange(_ctx, key, start, stop).Result()
}

// 从链表左侧弹出数据
func LPop(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.LPop(_ctx, key).Result()
}

// 从链表右侧弹出数据
func RPop(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.RPop(_ctx, key).Result()
}

// 根据值移除元素
func LRem(key string, count int64, value any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.LRem(_ctx, key, count, value).Result()
}
