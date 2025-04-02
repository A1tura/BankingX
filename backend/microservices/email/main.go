package main

import (
	"email/rabbitmq"
	"log"

	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
    if err := rabbitmq.Listen(); err != nil {
        log.Fatal(err)
    }
}
