package rateLimit

import (
	"errors"
	"github.com/fushiliang321/go-core/config/rateLimit"
	"sync"
	"time"
)

type Service struct{}

type tokenBucket struct {
	global int //全局令牌桶
	//paths       map[string]int //地址令牌桶
	lastUseTime int64 //最后使用的时间
	sync.Mutex
}

var (
	tokenBucketMap sync.Map
	configData     *rateLimit.RateLimit
)

func (Service) Start(_ *sync.WaitGroup) {
	configData = rateLimit.Get()
	if configData == nil {
		return
	}
	go func() {
		var (
			t      int64
			ok     bool
			bucket *tokenBucket
		)
		for {
			time.Sleep(time.Second)
			t = time.Now().Unix()
			tokenBucketMap.Range(func(key, value any) bool {
				bucket, ok = value.(*tokenBucket)
				if !ok || (t-bucket.lastUseTime) > configData.IdleTime {
					bucket = nil
					tokenBucketMap.Delete(key)
					return true
				}
				bucket.charge()
				return true
			})
		}
	}()
}

func Process(key string, path string) error {
	var bucket *tokenBucket
	v, ok := tokenBucketMap.Load(key)
	if !ok {
		bucket = bucketInit()
		tokenBucketMap.Store(key, bucket)
	} else {
		bucket = v.(*tokenBucket)
	}
	bucket.Lock()
	defer bucket.Unlock()
	if bucket.global < 1 {
		return errors.New("请求频率超出")
	}
	bucket.global = bucket.global - configData.Consume
	bucket.lastUseTime = time.Now().Unix()
	return nil
}

// 初始化令牌桶
func bucketInit() *tokenBucket {
	bucket := &tokenBucket{}
	bucket.global = configData.Capacity
	return bucket
}

// 令牌桶充能
func (bucket *tokenBucket) charge() {
	bucket.Lock()
	defer bucket.Unlock()
	bucket.global = bucket.global + configData.Create
	if bucket.global > configData.Capacity {
		bucket.global = configData.Capacity
	}
}
