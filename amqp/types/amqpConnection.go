package types

import (
	amqp3 "github.com/streadway/amqp"
	"sync"
)

type AmqpConnection struct {
	Producer *amqp3.Connection
	Consumer *amqp3.Connection
	sync.RWMutex
}
