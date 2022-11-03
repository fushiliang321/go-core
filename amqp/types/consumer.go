package types

type Result = byte
type ConsumerMessageHandle = func(data []byte) Result

type Consumer struct {
	Exchange            string
	RoutingKey          string
	Queue               string
	Type                string
	Durable             bool
	AutoDeletedExchange bool
	AutoDeletedQueue    bool
	Handler             ConsumerMessageHandle
}
