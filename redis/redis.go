package redis

import (
	"context"
	"fmt"
	redisConfig "github.com/fushiliang321/go-core/config/redis"
	"github.com/go-redis/redis/v9"
	"strconv"
)

var (
	_client *redis.Client
	_ctx    = context.Background()
)

func NewClient() *redis.Client {
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
		fmt.Println("connection redis error：", err.Error())
		return nil
	}
	return c
}

func client() *redis.Client {
	if _client == nil {
		_client = NewClient()
	}
	return _client
}
