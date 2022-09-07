package connection

import (
	config "github.com/fushiliang321/go-core/config/amqp"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type AmqpConnection struct {
	Producer *amqp.Connection
	Consumer *amqp.Connection
	sync.RWMutex
}

var connection = &AmqpConnection{}

func amqpInit() *AmqpConnection {
	connection.Lock()
	defer connection.Unlock()
	if amqpIsAvailable() {
		//可能获取到锁后已经有其他协程修改了数据
		return connection
	}
	config := config.Get()
	if config.Host == "" || config.Port == "" {
		return nil
	}
	var err error
	url := "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port
	consumer, err := amqp.Dial(url)
	if err != nil {
		log.Println("amqp err", err)
		return nil
	}
	producer, err := amqp.Dial(url)
	if err != nil {
		log.Println("amqp err", err)
		return nil
	}
	connection.Consumer = consumer
	connection.Producer = producer
	return connection
}

// 判断amqp是否可用
func amqpIsAvailable() bool {
	if connection.Consumer != nil && connection.Producer != nil {
		if !connection.Consumer.IsClosed() && !connection.Producer.IsClosed() {
			return true
		}
		connection.Consumer.Close()
		connection.Producer.Close()
	}
	return false
}

func GetAmqp() *AmqpConnection {
	if amqpIsAvailable() {
		return connection
	}
	return amqpInit()
}
