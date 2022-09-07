package connection

import (
	amqp2 "github.com/fushiliang321/go-core/config/amqp"
	amqp3 "github.com/streadway/amqp"
	"log"
	"sync"
)

type AmqpConnection struct {
	Producer *amqp3.Connection
	Consumer *amqp3.Connection
	sync.RWMutex
}

var amqp = &AmqpConnection{}

func amqpInit() *AmqpConnection {
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

func GetAmqp() *AmqpConnection {
	if amqpIsAvailable() {
		return amqp
	}
	return amqpInit()
}
