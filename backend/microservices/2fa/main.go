package main

import (
	"2fa/controllers"
	"2fa/service"
	"db"
	"log"
	"net/http"
	"os"

	"middlewares"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	db := db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"))

	redisCon := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	authRedisCon := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	rabbitmq, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := service.Start(rabbitmq, redisCon); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := middlewares.GetMiddleware(db, rabbitmq)
	ctx = middlewares.AddMiddleware(ctx, "redis", redisCon)
	ctx = middlewares.AddMiddleware(ctx, "redisAuth", authRedisCon)

	http.Handle("/verifyAuth", ctx(http.HandlerFunc(controllers.VerifyAuth)))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	select {}
}
