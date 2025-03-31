package middlewares

import (
	"context"
	"db"
	"net/http"
)

func GetMiddleware(db *db.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetContext(ctx context.Context) *db.DB {
	db := ctx.Value("db").(*db.DB)

	return db
}
