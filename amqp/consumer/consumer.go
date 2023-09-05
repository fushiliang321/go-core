package consumer

import (
	"fmt"
	"github.com/fushiliang321/go-core/amqp/connection"
	"github.com/fushiliang321/go-core/amqp/types"
	"github.com/fushiliang321/go-core/exception"
	"github.com/fushiliang321/go-core/helper/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type status = int8

type Consumer struct {
	*types.Consumer
	status  status
	channel *amqp.Channel
	sync.Mutex
}

const (
	ACK     = types.Result(0)
	NACK    = types.Result(1)
	REQUEUE = types.Result(2)
	REJECT  = types.Result(3)

	STATUS_START = 1
	STATUS_CLOSE = 0
)

func (consumer *Consumer) Start() {
	consumer.status = STATUS_START
	consumer.monitor()
}

func (consumer *Consumer) monitor() {
	consumer.Lock()
	defer consumer.Unlock()
	if consumer.status != STATUS_START {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error(consumer.Exchange+" amqp monitor recover:", fmt.Sprint(err))
			exception.Listener(consumer.Exchange+" amqp monitor", err)
			// 监听异常 要重试
			go consumer.retryMonitor()
		}
	}()
	_amqp := connection.GetAmqp()
	if _amqp == nil {
		// 监听失败 要重试
		go consumer.retryMonitor()
		return
	}
	var err error
	consumer.channel, err = _amqp.Consumer.Channel()
	if err != nil {
		logger.Warn(consumer.Exchange, "consumer channel error", err)
		// 监听失败 要重试
		go consumer.retryMonitor()
		return
	}
	closeChannel := true
	defer func() {
		if closeChannel {
			if consumer.channel != nil {
				consumer.channel.Close()
			}
			// 监听失败 要重试
			go consumer.retryMonitor()
		}
	}()
	if err != nil {
		logger.Warn(consumer.Exchange, "consumer channel error", err)
		return
	}
	err = channelInit(consumer.channel, consumer.Exchange, consumer.RoutingKey, consumer.Queue, consumer.Type, consumer.Durable, consumer.AutoDeletedExchange, consumer.AutoDeletedQueue)
	if err != nil {
		return
	}
	msgs, err := consumer.channel.Consume(consumer.Queue, "", false, false, false, true, amqp.Table{})
	if err != nil {
		logger.Warn(consumer.Exchange, "consumer consume error", err)
		return
	}
	closeChannel = false
	go func() {
		fun := func(d *amqp.Delivery) {
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
		logger.Info(consumer.Exchange, "channel close")
		// 断开后 要重新监听
		go consumer.retryMonitor()
	}()
}

// 重新监听
func (consumer *Consumer) retryMonitor() {
	if consumer.status != STATUS_START {
		return
	}
	time.Sleep(time.Second * 5)
	consumer.monitor()
}

// 关闭监听
func (consumer *Consumer) Close() {
	consumer.Lock()
	defer consumer.Unlock()
	consumer.status = STATUS_CLOSE
	if consumer.channel != nil {
		consumer.channel.Close()
		consumer.channel = nil
	}
}

func channelInit(channel *amqp.Channel, Exchange string, RoutingKey string, Queue string, kind string, durable bool, AutoDeletedExchange bool, AutoDeletedQueue bool) (err error) {
	if kind == "" {
		kind = types.ExchangeTypeDirect
	}
	err = channel.ExchangeDeclare(Exchange, kind, durable, AutoDeletedExchange, false, false, nil)
	if err != nil {
		logger.Warn("consumer exchange error", err)
		return
	}
	if Queue == "" {
		Queue = RoutingKey
	}
	q, err := channel.QueueDeclare(
		Queue,
		durable,
		AutoDeletedQueue,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Warn("consumer queue error", err)
		return
	}
	err = channel.QueueBind(q.Name, RoutingKey, Exchange, false, nil)
	if err != nil {
		logger.Warn("consumer queue error", err)
		return
	}
	return
}
