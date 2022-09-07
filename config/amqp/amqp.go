package amqp

import (
	"github.com/fushiliang321/go-core/amqp/types"
)

type Amqp struct {
	Host      string
	Port      string
	User      string
	Password  string
	Consumers []*types.Consumer
}

var amqp = &Amqp{
	Consumers: []*types.Consumer{},
}

func Set(config *Amqp) {
	amqp = config
}

func Get() *Amqp {
	return amqp
}
