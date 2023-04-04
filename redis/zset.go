package redis

import "github.com/go-redis/redis/v9"

type (
	Z        = redis.Z
	ZRangeBy = redis.ZRangeBy
)

// 添加元素
func ZAdd(key string, members ...Z) (int64, error) {
	return client().ZAdd(_ctx, key, members...).Result()
}

// 增加元素分值
func ZIncrBy(key string, increment float64, member string) (float64, error) {
	return client().ZIncrBy(_ctx, key, increment, member).Result()
}

// 获取根据score排序后的数据段，升序
func ZRange(key string, startStop ...int64) ([]string, error) {
	var start, stop int64
	switch len(startStop) {
	case 0:
		start = 0
		stop = -1
	case 1:
		start = startStop[0]
		stop = -1
	default:
		start = startStop[0]
		stop = startStop[1]
	}
	return client().ZRange(_ctx, key, start, stop).Result()
}

// 获取根据score排序后的数据段，降序
func ZRevRange(key string, startStop ...int64) ([]string, error) {
	var start, stop int64
	switch len(startStop) {
	case 0:
		start = 0
		stop = -1
	case 1:
		start = startStop[0]
		stop = -1
	default:
		start = startStop[0]
		stop = startStop[1]
	}
	return client().ZRevRange(_ctx, key, start, stop).Result()
}

// 获取score过滤后排序的数据段，升序
func ZRangeByScore(key string, opt *ZRangeBy) ([]string, error) {
	return client().ZRangeByScore(_ctx, key, opt).Result()
}

// 获取score过滤后排序的数据段，降序
func ZRevRangeByScore(key string, opt *ZRangeBy) ([]string, error) {
	return client().ZRevRangeByScore(_ctx, key, opt).Result()
}

// 获取元素个数
func ZCard(key string) (int64, error) {
	return client().ZCard(_ctx, key).Result()
}

// 获取区间内元素个数
func ZCount(key, min, max string) (int64, error) {
	return client().ZCount(_ctx, key, min, max).Result()
}

// 获取元素的score
func ZScore(key, member string) (float64, error) {
	return client().ZScore(_ctx, key, member).Result()
}

// 获取某个元素在集合中的排名，升序
func ZRank(key, member string) (int64, error) {
	return client().ZRank(_ctx, key, member).Result()
}

// 获取某个元素在集合中的排名，降序
func ZRevRank(key, member string) (int64, error) {
	return client().ZRevRank(_ctx, key, member).Result()
}

// 删除元素
func ZRem(key string, members ...any) (int64, error) {
	return client().ZRem(_ctx, key, members...).Result()
}

// 根据排名来删除
func ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return client().ZRemRangeByRank(_ctx, key, start, stop).Result()
}

// 根据分值区间来删除
func ZRemRangeByScore(key, start, stop string) (int64, error) {
	return client().ZRemRangeByScore(_ctx, key, start, stop).Result()
}
