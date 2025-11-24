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

	gameState := gamelogic.NewGameState(user)

	for {
		cmd := gamelogic.GetInput()
		switch cmd[0] {
		case "spawn":
			err = gameState.CommandSpawn(cmd)
			if err != nil {
				log.Fatal(err)
			}
		case "move":
			_, err := gameState.CommandMove(cmd)
			if err != nil {
				log.Fatal(err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
			continue
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Print("unknown command")
			continue
		}
	}
}
