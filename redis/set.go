package redis

// 添加元素
func SAdd(key string, members ...any) (int64, error) {
	return client().SAdd(_ctx, key, members...).Result()
}

// 随机获取一个元素
func SRandMember(key string) (string, error) {
	return client().SRandMember(_ctx, key).Result()
}

// 随机获取多个元素
func SRandMemberN(key string, count int64) ([]string, error) {
	return client().SRandMemberN(_ctx, key, count).Result()
}

// 随机移除一个元素,并返回
func SPop(key string) (string, error) {
	return client().SPop(_ctx, key).Result()
}

// 随机移除多个元素,并返回
func SPopN(key string, count int64) ([]string, error) {
	return client().SPopN(_ctx, key, count).Result()
}

// 删除集合里指定的值
func SRem(key string, members ...any) (int64, error) {
	return client().SRem(_ctx, key, members...).Result()
}

// 获取所有成员
func SMembers(key string) ([]string, error) {
	return client().SMembers(_ctx, key).Result()
}

// 判断元素是否在集合中
func SIsMember(key string, member any) (bool, error) {
	return client().SIsMember(_ctx, key, member).Result()
}

// 获取集合元素个数
func SCard(key string) (int64, error) {
	return client().SCard(_ctx, key).Result()
}

// 并集
func SUnion(keys ...string) ([]string, error) {
	return client().SUnion(_ctx, keys...).Result()
}

// 差集
func SDiff(keys ...string) ([]string, error) {
	return client().SDiff(_ctx, keys...).Result()
}

// 交集
func SInter(keys ...string) ([]string, error) {
	return client().SInter(_ctx, keys...).Result()
}
