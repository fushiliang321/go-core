package amqp

import (
	"github.com/fushiliang321/go-core/amqp/connection"
	"github.com/fushiliang321/go-core/amqp/consumer"
	"github.com/fushiliang321/go-core/amqp/types"
	amqp2 "github.com/fushiliang321/go-core/config/amqp"
	"github.com/fushiliang321/go-core/helper"
	amqp3 "github.com/streadway/amqp"
	"log"
	"sync"
)

type Service struct{}

func (Service) Start(_ *sync.WaitGroup) {
	config := amqp2.Get()
	if len(config.Consumers) > 0 {
		//有消费者
		for _, _consumer := range config.Consumers {
			con := consumer.Consumer{
				Consumer: _consumer,
			}
			con.Start()
		}
	}
}

func Publish(producer *types.Producer) {
	Amqp := connection.GetAmqp()
	if Amqp == nil {
		return
	}
	var err error
	channel, err := Amqp.Producer.Channel()
	if err != nil {
		log.Println("producer channel error", err)
		return
	}
	defer channel.Close()
	if err != nil {
		return
	}
	body, err := helper.AnyToBytes(producer.Data)
	deliveryMode := amqp3.Persistent
	if !producer.Persistence {
		deliveryMode = amqp3.Transient
	}
	err = channel.Publish(producer.Exchange, producer.RoutingKey, false, false, amqp3.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: deliveryMode,
		Body:         body,
		Expiration:   producer.Expiration,
		Priority:     producer.Priority,
	})
	if err != nil {
		log.Println("publish producer error", err)
	}
}
