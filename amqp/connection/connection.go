package connection

import (
	config "github.com/fushiliang321/go-core/config/amqp"
	"github.com/fushiliang321/go-core/helper/logger"
	amqp "github.com/rabbitmq/amqp091-go"
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
	configData := config.Get()
	if configData.Host == "" || configData.Port == "" {
		return nil
	}
	var err error
	url := "amqp://" + configData.User + ":" + configData.Password + "@" + configData.Host + ":" + configData.Port
	consumer, err := amqp.Dial(url)
	if err != nil {
		logger.Warn("amqp err", err)
		return nil
	}
	producer, err := amqp.Dial(url)
	if err != nil {
		logger.Warn("amqp err", err)
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
