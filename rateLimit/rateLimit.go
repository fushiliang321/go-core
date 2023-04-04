package rateLimit

import (
	"github.com/fushiliang321/go-core/config/rateLimit"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Service     struct{}
	tokenBucket struct {
		global atomic.Int32
		//paths       map[string]int //地址令牌桶
		lastUseTime int64 //最后使用的时间
		sync.Mutex
	}
)

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
			bucket = nil
		}
	}()
}

func Process(key string, path string) bool {
	var bucket *tokenBucket
	if v, ok := tokenBucketMap.Load(key); ok {
		bucket = v.(*tokenBucket)
	} else {
		bucket = bucketInit()
		tokenBucketMap.Store(key, bucket)
	}
	bucket.Lock()
	defer bucket.Unlock()
	if bucket.global.Load() < 1 {
		return false
	}
	bucket.global.Add(-configData.Consume)
	bucket.lastUseTime = time.Now().Unix()
	return true
}

// 初始化令牌桶
func bucketInit() *tokenBucket {
	bucket := &tokenBucket{}
	bucket.global.Store(configData.Capacity)
	return bucket
}

// 令牌桶充能
func (bucket *tokenBucket) charge() {
	if bucket.global.Add(configData.Create) > configData.Capacity {
		bucket.global.Store(configData.Capacity)
	}
}
