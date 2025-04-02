package middlewares

import (
	"context"
	"db"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Services struct {
	DB       *db.DB
	Rabbitmq *amqp.Connection
}

type AuthInfo struct {
	IsAuth bool
	UserId int
}

func GetMiddleware(db *db.DB, rabbitmq *amqp.Connection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			ctx = context.WithValue(ctx, "rabbitmq", rabbitmq)

			authHeader := strings.Split(r.Header.Get("Authorization"), " ")

			if authHeader[0] == "Bearer" {
				token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					ctx = context.WithValue(ctx, "isAuth", false)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				if !token.Valid {
					ctx = context.WithValue(ctx, "isAuth", false)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					ctx = context.WithValue(ctx, "isAuth", false)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				userId := claims["userId"].(float64)

				ctx = context.WithValue(ctx, "isAuth", true)
				ctx = context.WithValue(ctx, "userId", int(userId))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuth(ctx context.Context) *AuthInfo {
	isAuthValue := ctx.Value("isAuth")
	userIdValue := ctx.Value("userId")

	if isAuthValue == nil || userIdValue == nil {
		return &AuthInfo{IsAuth: false, UserId: 0}
	}

    isAuth := isAuthValue.(bool)
    userId := userIdValue.(int)

	authInfo := AuthInfo{
		IsAuth: isAuth,
		UserId: userId,
	}

	return &authInfo
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
