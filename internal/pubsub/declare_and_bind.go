package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

type QueueSettings struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
}

func (s SimpleQueueType) Settings() QueueSettings {
	switch s {
	case Durable:
		return QueueSettings{
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
		}
	case Transient:
		return QueueSettings{
			Durable:    false,
			AutoDelete: true,
			Exclusive:  true,
		}
	}

	return QueueSettings{}
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	settings := queueType.Settings()

	queue, err := ch.QueueDeclare(queueName, settings.Durable, settings.AutoDelete, settings.Exclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
