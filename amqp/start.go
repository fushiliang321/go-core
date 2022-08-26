package amqp

import (
	"encoding/json"
	"fmt"
	"gitee.com/zvc/go-core/amqp/types"
	amqp2 "gitee.com/zvc/go-core/config/amqp"
	"gitee.com/zvc/go-core/exception"
	amqp3 "github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type Service struct {
}

func (Service) Start(wg *sync.WaitGroup) {
	config := amqp2.Get()
	if len(config.Consumers) > 0 {
		//有消费者
		for i := range config.Consumers {
			con := config.Consumers[i]
			monitor(con)
		}
	}
}

func monitor(consumer *types.Consumer) {
	defer func() {
		if err := recover(); err != nil {
			exception.Listener("amqp monitor", err)

			// 监听异常 要重试
			go func(c *types.Consumer) {
				time.Sleep(time.Second * 5)
				monitor(c)
			}(consumer)
		}
	}()
	if getAmqp() == nil {
		// 监听失败 要重试
		go func(c *types.Consumer) {
			time.Sleep(time.Second * 5)
			monitor(c)
		}(consumer)
		return
	}
	channel, err := getAmqp().Consumer.Channel()
	closeChannel := true
	defer func() {
		if closeChannel {
			channel.Close()

			// 监听失败 要重试
			go func(c *types.Consumer) {
				time.Sleep(time.Second * 5)
				monitor(c)
			}(consumer)
		}
	}()
	if err != nil {
		log.Println("consumer channel error", err)
		return
	}
	err = channelInit(channel, consumer.Exchange, consumer.RoutingKey, consumer.Queue, consumer.Type, consumer.Durable)
	if err != nil {
		return
	}
	msgs, err := channel.Consume(consumer.Queue, "", false, false, false, true, amqp3.Table{})
	if err != nil {
		log.Println("consumer consume error", err)
		return
	}
	closeChannel = false
	handler := consumer.Handler
	go func() {
		for d := range msgs {
			handler(d.Body, d)
		}
		log.Println("channel close")
		// 断开后 要重新监听
		go func(c *types.Consumer) {
			time.Sleep(time.Second * 5)
			monitor(c)
		}(consumer)
	}()
}

func channelInit(channel *amqp3.Channel, Exchange string, RoutingKey string, Queue string, kind string, durable bool) (err error) {
	if kind == "" {
		kind = ExchangeTypeDirect
	}
	err = channel.ExchangeDeclare(Exchange, kind, durable, false, false, false, nil)
	if err != nil {
		log.Println("consumer exchange error", err)
		return
	}
	if Queue == "" {
		Queue = RoutingKey
	}
	q, err := channel.QueueDeclare(
		Queue,
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("consumer queue error", err)
		return
	}
	err = channel.QueueBind(q.Name, RoutingKey, Exchange, false, nil)
	if err != nil {
		log.Println("consumer queue error", err)
		return
	}
	return
}

func Publish(producer *types.Producer) {
	Amqp := getAmqp()
	if Amqp == nil {
		return
	}
	channel, err := Amqp.Producer.Channel()
	if err != nil {
		log.Println("producer channel error", err)
		return
	}
	defer channel.Close()
	if err != nil {
		return
	}
	marshal, err := json.Marshal(producer.Data)
	if err != nil {
		return
	}
	err = channel.Publish(producer.Exchange, producer.RoutingKey, false, false, amqp3.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp3.Persistent,
		Body:         marshal,
	})
	if err != nil {
		log.Println("publish producer error", err)
	}
	fmt.Println("Publish  end")
}
