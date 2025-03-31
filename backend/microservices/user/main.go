package main

import (
	"db"
	"middlewares"
	"net/http"

	"user/controllers"
)

var con = db.Connect("admin", "admin", "bank")
var ctxMiddleware = middlewares.GetMiddleware(con)

func main() {
	http.Handle("/signup", ctxMiddleware(http.HandlerFunc(controllers.SignUp)))

	http.ListenAndServe(":8080", nil)
}
