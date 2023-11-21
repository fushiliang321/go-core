package redis

// 添加元素
func SAdd(key string, members ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.SAdd(_ctx, key, members...).Result()
}

// 随机获取一个元素
func SRandMember(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.SRandMember(_ctx, key).Result()
}

// 随机获取多个元素
func SRandMemberN(key string, count int64) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SRandMemberN(_ctx, key, count).Result()
}

// 随机移除一个元素,并返回
func SPop(key string) (string, error) {
	c, err := client()
	if err != nil {
		return "", err
	}
	return c.SPop(_ctx, key).Result()
}

// 随机移除多个元素,并返回
func SPopN(key string, count int64) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SPopN(_ctx, key, count).Result()
}

// 删除集合里指定的值
func SRem(key string, members ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.SRem(_ctx, key, members...).Result()
}

// 获取所有成员
func SMembers(key string) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SMembers(_ctx, key).Result()
}

// 判断元素是否在集合中
func SIsMember(key string, member any) (bool, error) {
	c, err := client()
	if err != nil {
		return false, err
	}
	return c.SIsMember(_ctx, key, member).Result()
}

// 获取集合元素个数
func SCard(key string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.SCard(_ctx, key).Result()
}

// 并集
func SUnion(keys ...string) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SUnion(_ctx, keys...).Result()
}

// 差集
func SDiff(keys ...string) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SDiff(_ctx, keys...).Result()
}

// 交集
func SInter(keys ...string) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.SInter(_ctx, keys...).Result()
}
