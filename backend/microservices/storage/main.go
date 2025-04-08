package main

import (
	"log"
	"net/http"
	"os"
	"storage/controllers"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"

	"db"
	"middlewares"
)

func main() {
	godotenv.Load()

	con := db.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"))
	rabbitmq, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatalf("Error while connecting to rabbitmq: ", err)
	}

	minio, err := minio.New(os.Getenv("MINIO_HOST"), &minio.Options{
		Creds: credentials.NewStaticV4(os.Getenv("MINIO_USERNAME"), os.Getenv("MINIO_PASSWORD"), ""),
	})

	var ctxMiddleware = middlewares.GetMiddleware(con, rabbitmq)
	ctxMiddleware = middlewares.AddMiddleware(ctxMiddleware, "minio", minio)

	// KYC storage
	http.Handle("/KYC/document", ctxMiddleware(http.HandlerFunc(controllers.Document)))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
