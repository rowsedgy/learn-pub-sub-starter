package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("error connecting to rabbitmq instance: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// new durable queue
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", "game_logs.*", pubsub.Durable)
	if err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintServerHelp()
	fmt.Println("Connection successful!")
	fmt.Println("Choose a command:")

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			log.Print("Sending a pause message...")
			if err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect,
				string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: true,
				}); err != nil {
				log.Fatalf("error publishing pause message: %v", err)
			}
		case "resume":
			log.Print("Sending resume message...")
			if err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect,
				string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: false,
				}); err != nil {
				log.Fatalf("error publishing resume message: %v", err)
			}
		case "quit":
			log.Print("Quitting...")
			return
		default:
			log.Print("Command not valid")
		}
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan

	// fmt.Println("Received signal")

}
