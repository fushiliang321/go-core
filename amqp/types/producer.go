package types

type Producer struct {
	Exchange    string
	RoutingKey  string
	Data        any
	Expiration  string //过期时间，毫秒
	Persistence bool   //持久化消息
	Priority    uint8  //优先级0-9
}
