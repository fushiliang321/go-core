package amqp

import (
	amqp2 "github.com/fushiliang321/go-core/amqp/consumer"
)

type Amqp struct {
	Host      string
	Port      string
	User      string
	Password  string
	Consumers []*amqp2.Consumer
}

var amqp = &Amqp{
	Consumers: []*amqp2.Consumer{},
}

func Set(config *Amqp) {
	amqp = config
}

func Get() *Amqp {
	return amqp
}
