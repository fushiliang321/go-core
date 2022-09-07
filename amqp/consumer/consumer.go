package consumer

import (
	"github.com/fushiliang321/go-core/amqp/connection"
	"github.com/fushiliang321/go-core/amqp/types"
	"github.com/fushiliang321/go-core/exception"
	amqp3 "github.com/streadway/amqp"
	"log"
	"time"
)

type Consumer struct {
	*types.Consumer
}

const (
	ACK     = types.Result(0)
	NACK    = types.Result(1)
	REQUEUE = types.Result(2)
	REJECT  = types.Result(3)
)

func (consumer *Consumer) Monitor() {
	defer func() {
		if err := recover(); err != nil {
			exception.Listener("amqp monitor", err)
			// 监听异常 要重试
			go consumer.retryMonitor()
		}
	}()
	amqp := connection.GetAmqp()
	if amqp == nil {
		// 监听失败 要重试
		go consumer.retryMonitor()
		return
	}
	channel, err := amqp.Consumer.Channel()
	if err != nil {
		log.Println("consumer channel error", err)
		// 监听失败 要重试
		go consumer.retryMonitor()
		return
	}
	closeChannel := true
	defer func() {
		if closeChannel {
			channel.Close()
			// 监听失败 要重试
			go consumer.retryMonitor()
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
	go func() {
		fun := func(d *amqp3.Delivery) {
			switch consumer.Handler(d.Body) {
			case ACK:
				d.Ack(false)
			case NACK:
				d.Nack(false, false)
			case REQUEUE:
				d.Reject(true)
			case REJECT:
				d.Reject(false)
			}
		}
		for d := range msgs {
			fun(&d)
		}
		log.Println("channel close")
		// 断开后 要重新监听
		go consumer.retryMonitor()
	}()
}

func (consumer *Consumer) retryMonitor() {
	time.Sleep(time.Second * 5)
	consumer.Monitor()
}

func channelInit(channel *amqp3.Channel, Exchange string, RoutingKey string, Queue string, kind string, durable bool) (err error) {
	if kind == "" {
		kind = types.ExchangeTypeDirect
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
