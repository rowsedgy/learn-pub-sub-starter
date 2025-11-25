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

	// pause subscribe
	if err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gameState)); err != nil {
		log.Fatal(err)
	}

	// move subscribe
	moveQueueName := fmt.Sprintf("army_moves.%s", user)
	if err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, moveQueueName, routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerMove(gameState)); err != nil {
		log.Fatal(err)
	}

	// move channel
	moveCh, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, moveQueueName, moveQueueName, pubsub.Transient)
	if err != nil {
		log.Fatal(err)
	}

	for {
		cmd := gamelogic.GetInput()
		switch cmd[0] {
		case "spawn":
			err = gameState.CommandSpawn(cmd)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			move, err := gameState.CommandMove(cmd)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if err = pubsub.PublishJSON(moveCh, routing.ExchangePerilTopic, moveQueueName, move); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Move published succesfully")
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
