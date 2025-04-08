package middlewares

import (
	"context"
	"db"
	"middlewares/dal"
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
	IsAuth         bool
	UserId         int
	EmailConfirmed *bool
	KYCStatus      *string
}

func GetMiddleware(db *db.DB, rabbitmq *amqp.Connection) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			ctx = context.WithValue(ctx, "rabbitmq", rabbitmq)
			ctx = context.WithValue(ctx, "authInfo", &AuthInfo{
				IsAuth:         false,
				UserId:         0,
				EmailConfirmed: nil,
				KYCStatus:      nil,
			})

			authHeader := strings.Split(r.Header.Get("Authorization"), " ")

			if len(authHeader) == 2 {
				if authHeader[0] == "Bearer" {
					token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
						return []byte(os.Getenv("JWT_SECRET")), nil
					})

					if err != nil {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}

					if !token.Valid {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}

					claims, ok := token.Claims.(jwt.MapClaims)
					if !ok {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					userId := claims["userId"].(float64)

					authInfo := AuthInfo{
						IsAuth: true,
						UserId: int(userId),
					}

					emailConfirmed, err := dal.EmailConfirmed(db, int(userId))
					if err != nil {
						authInfo.EmailConfirmed = nil
					} else {
						authInfo.EmailConfirmed = &emailConfirmed
					}

					KYCStatus, err := dal.KYCStatus(db, int(userId))
					if err != nil {
						authInfo.KYCStatus = nil
					} else {
						authInfo.KYCStatus = &KYCStatus
					}

					ctx = context.WithValue(ctx, "authInfo", &authInfo)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AddMiddleware(f func(next http.Handler) http.Handler, name string, value any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return f(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), name, value)

			next.ServeHTTP(w, r.WithContext(ctx))
		}))
	}
}

func GetAuth(ctx context.Context) *AuthInfo {
	authInfo := ctx.Value("authInfo").(*AuthInfo)

	return authInfo
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
