package redis

import (
	"context"
	redisConfig "github.com/fushiliang321/go-core/config/redis"
	"github.com/fushiliang321/go-core/event/handles/core"
	"github.com/fushiliang321/go-core/helper/logger"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

var (
	_client *redis.Client
	_lock   sync.RWMutex
	_ctx    = context.Background()
)

func NewClient() *redis.Client {
	core.AwaitStartFinish()
	var (
		config = redisConfig.Get()
		c      = redis.NewClient(&redis.Options{
			Addr:     config.Host + ":" + strconv.Itoa(config.Port),
			Password: config.Password,
			DB:       config.Db,
		})
		_, err = c.Ping(_ctx).Result()
	)
	if err != nil {
		logger.Warn("connection redis error：" + err.Error())
		return nil
	}
	return c
}

func client() *redis.Client {
	if _client == nil {
		_lock.Lock()
		if _client == nil {
			_client = NewClient()
		}
		_lock.Unlock()
	}
	return _client
}
