package main

import (
	"db"
	"middlewares"
	"net/http"

	"github.com/joho/godotenv"

	"user/controllers"
)

var con = db.Connect("admin", "admin", "bank")
var ctxMiddleware = middlewares.GetMiddleware(con)

func main() {
	godotenv.Load()
	http.Handle("/signup", ctxMiddleware(http.HandlerFunc(controllers.SignUp)))
	http.Handle("/signin", ctxMiddleware(http.HandlerFunc(controllers.SignIn)))

	http.ListenAndServe(":8080", nil)
}
