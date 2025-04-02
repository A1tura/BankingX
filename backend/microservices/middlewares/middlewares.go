package middlewares

import (
	"context"
	"db"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

type Services struct {
	DB       *db.DB
	Rabbitmq *amqp.Connection
}

func GetMiddleware(db *db.DB, rabbitmq *amqp.Connection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			ctx = context.WithValue(ctx, "rabbitmq", rabbitmq)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetContext(ctx context.Context) *Services {
	db := ctx.Value("db").(*db.DB)
	rabbtimq := ctx.Value("rabbitmq").(*amqp.Connection)

	services := Services{
		DB:       db,
		Rabbitmq: rabbtimq,
	}

	return &services
}
