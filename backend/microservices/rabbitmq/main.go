package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

    return err
}

func main() {
	godotenv.Load()
	con, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatal(err)
	}

    queues := []string{
        "email",
    }

	ch, err := con.Channel()
	if err != nil {
		log.Fatal(err)
	}

    for _, queue := range queues {
        declareQueue(ch, queue)
    }
}
