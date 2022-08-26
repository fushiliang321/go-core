package types

import (
	amqp2 "github.com/streadway/amqp"
)

type ConsumerMessageHandle = func(data []byte, delivery amqp2.Delivery)

type Consumer struct {
	Exchange   string
	RoutingKey string
	Queue      string
	Type       string
	Durable    bool
	Handler    ConsumerMessageHandle
}
