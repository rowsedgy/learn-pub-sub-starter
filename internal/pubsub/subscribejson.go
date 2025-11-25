package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackReque
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	consumeCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return (err)
	}

	go func() {
		for v := range consumeCh {
			var receiver T
			if err = json.Unmarshal(v.Body, &receiver); err != nil {
				log.Fatal(err)
			}
			ackType := handler(receiver)
			switch ackType {
			case Ack:
				_ = v.Ack(false)
				fmt.Println("Ack sent")
			case NackReque:
				_ = v.Nack(false, true)
				fmt.Println("NackRequeue sent")
			case NackDiscard:
				_ = v.Nack(false, false)
				fmt.Println("NackDiscard sent")
			}

		}
	}()

	return nil
}
