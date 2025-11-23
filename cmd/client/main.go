package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatal(err)
	}

	user, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error welcoming user: %v", err)
	}

	queueName := fmt.Sprintf("pause.%s", user)
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, queueName, pubsub.Transient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Queue created!")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

}
