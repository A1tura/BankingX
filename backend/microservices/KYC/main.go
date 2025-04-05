package main

import (
	"db"
	"kyc/controllers"
	"log"
	"middlewares"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load()

	db := db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"))
	rabbitmq, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatalf("Error while connecting to rabbitmq: ", rabbitmq)
	}

	middleware := middlewares.GetMiddleware(db, rabbitmq)

	http.Handle("/test", middleware(http.HandlerFunc(controllers.Test)))
	http.Handle("/", middleware(http.HandlerFunc(controllers.KYC)))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
