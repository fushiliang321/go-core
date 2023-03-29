package rateLimit

type RateLimit struct {
	Create   int   //每秒生成令牌数
	Consume  int   //每次请求消耗令牌数
	Capacity int   //令牌桶最大容量
	IdleTime int64 //闲置时长（s）
}

var config = &RateLimit{}

func Set(c *RateLimit) {
	config = c
}

func Get() *RateLimit {
	return config
}
