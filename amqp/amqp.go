package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/fushiliang321/go-core/amqp/types"
	amqp2 "github.com/fushiliang321/go-core/config/amqp"
	amqp3 "github.com/streadway/amqp"
	"log"
)

const (
	ExchangeTypeDirect  = "direct"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeHeaders = "headers"
)

var amqp = &types.AmqpConnection{}

func amqpInit() *types.AmqpConnection {
	amqp.Lock()
	defer amqp.Unlock()
	if amqpIsAvailable() {
		//可能获取到锁后已经有其他协程修改了数据
		return amqp
	}
	config := amqp2.Get()
	if config.Host == "" || config.Port == "" {
		return nil
	}
	var err error
	url := "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port
	consumer, err := amqp3.Dial(url)
	if err != nil {
		log.Println("amqp err", err)
		return nil
	}
	producer, err := amqp3.Dial(url)
	if err != nil {
		log.Println("amqp err", err)
		return nil
	}
	amqp.Consumer = consumer
	amqp.Producer = producer
	return amqp
}

// 判断amqp是否可用
func amqpIsAvailable() bool {
	if amqp.Consumer != nil && amqp.Producer != nil {
		if !amqp.Consumer.IsClosed() && !amqp.Producer.IsClosed() {
			return true
		}
		amqp.Consumer.Close()
		amqp.Producer.Close()
	}
	return false
}

func GetAmqp() *types.AmqpConnection {
	if amqpIsAvailable() {
		return amqp
	}
	return amqpInit()
}

func Publish(producer *types.Producer) {
	Amqp := GetAmqp()
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
