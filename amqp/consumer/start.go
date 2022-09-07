package consumer

import (
	amqp2 "github.com/fushiliang321/go-core/config/amqp"
	"sync"
)

type Service struct {
}

func (Service) Start(wg *sync.WaitGroup) {
	config := amqp2.Get()
	if len(config.Consumers) > 0 {
		//有消费者
		for i := range config.Consumers {
			con := config.Consumers[i]
			con.Monitor()
		}
	}
}
