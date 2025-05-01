package main

import (
	"db"
	"log"
	"middlewares"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"

	"accounts/controllers"
)

func main() {
	godotenv.Load()

	con := db.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"))
	rabbitmq, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatalf("Error while connecting to rabbitmq: ", err)
	}
	var ctxMiddleware = middlewares.GetMiddleware(con, rabbitmq)

	http.Handle("/accounts", ctxMiddleware(http.HandlerFunc(controllers.Accounts)))
	http.Handle("/createAccount", ctxMiddleware(http.HandlerFunc(controllers.CreateAccount)))

	http.Handle("/frozeAccount", ctxMiddleware(http.HandlerFunc(controllers.FrozeAccount)))
	http.Handle("/unfrozeAccount", ctxMiddleware(http.HandlerFunc(controllers.UnfrozeAccount)))
	// 	http.Handle("/closeAccount", ctxMiddleware(http.HandlerFunc(controllers.CloseAccount)))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
