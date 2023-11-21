package redis

import "github.com/redis/go-redis/v9"

type (
	Z        = redis.Z
	ZRangeBy = redis.ZRangeBy
)

// ZAdd 添加元素
func ZAdd(key string, members ...Z) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZAdd(_ctx, key, members...).Result()
}

// ZIncrBy 增加元素分值
func ZIncrBy(key string, increment float64, member string) (float64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZIncrBy(_ctx, key, increment, member).Result()
}

// ZRange 获取根据score排序后的数据段，升序
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
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.ZRange(_ctx, key, start, stop).Result()
}

// ZRevRange 获取根据score排序后的数据段，降序
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
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.ZRevRange(_ctx, key, start, stop).Result()
}

// ZRangeByScore 获取score过滤后排序的数据段，升序
func ZRangeByScore(key string, opt *ZRangeBy) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.ZRangeByScore(_ctx, key, opt).Result()
}

// ZRevRangeByScore 获取score过滤后排序的数据段，降序
func ZRevRangeByScore(key string, opt *ZRangeBy) ([]string, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	return c.ZRevRangeByScore(_ctx, key, opt).Result()
}

// ZCard 获取元素个数
func ZCard(key string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZCard(_ctx, key).Result()
}

// ZCount 获取区间内元素个数
func ZCount(key, min, max string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZCount(_ctx, key, min, max).Result()
}

// ZScore 获取元素的score
func ZScore(key, member string) (float64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZScore(_ctx, key, member).Result()
}

// ZRank 获取某个元素在集合中的排名，升序
func ZRank(key, member string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZRank(_ctx, key, member).Result()
}

// ZRevRank 获取某个元素在集合中的排名，降序
func ZRevRank(key, member string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZRevRank(_ctx, key, member).Result()
}

// ZRem 删除元素
func ZRem(key string, members ...any) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZRem(_ctx, key, members...).Result()
}

// ZRemRangeByRank 根据排名来删除
func ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZRemRangeByRank(_ctx, key, start, stop).Result()
}

// ZRemRangeByScore 根据分值区间来删除
func ZRemRangeByScore(key, start, stop string) (int64, error) {
	c, err := client()
	if err != nil {
		return 0, err
	}
	return c.ZRemRangeByScore(_ctx, key, start, stop).Result()
}
